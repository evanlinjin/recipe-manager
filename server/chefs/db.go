package chefs

import (
	"sync"

	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"net/smtp"
	"strings"
	"time"
)

type ChefsDB struct {
	sync.Mutex
	config  *Config
	session *mgo.Session
	chefs   *mgo.Collection
	verts   *mgo.Collection
}

func MakeChefsDB(config *Config) (c ChefsDB, e error) {
	if config == nil {
		e = fmt.Errorf("nil config")
		return
	}
	c.config = config
	c.session, e = mgo.Dial(config.MongoUrls)
	if e != nil {
		return
	}
	c.chefs = c.session.DB(DBAuth).C(CTChefs)
	c.verts = c.session.DB(DBAuth).C(CTVerts)
	return
}

func (c *ChefsDB) Close() {
	c.Lock()
	defer c.Unlock()
	c.session.Close()
}

func (c *ChefsDB) AddChef(email, pwd string) error {
	c.Lock()
	defer c.Unlock()
	email = strings.TrimSpace(email)

	// See if it's the first user. If so, make admin.
	// Else, check if this email is allowed.
	admin := false
	if n, _ := c.chefs.Count(); n == 0 {
		admin = true
	} else {
		e := c.session.DB(DBAuth).C(CTAllowed).Find(bson.M{"email": email})
		if e != nil {
			return ErrPermissionDenied
		}
	}

	// Check if chef already exists with specified email.
	temp := Chef{}
	if c.chefs.Find(bson.M{"email": email}).One(&temp) == nil {
		return &ErrChefAlreadyExists{email}
	}

	// Generate chef id, pwd_salt and pwd_hash.
	id := bson.NewObjectId() // id is string.
	pwdSalt := GetRand32Str()
	pwdHash, e := bcrypt.GenerateFromPassword([]byte(pwd+pwdSalt), 10)
	if e != nil {
		return &ErrInternal{e}
	}

	// Generate account activation key.
	actKey := GetRand32Str() + GetRand32Str()
	actKeySalt := GetRand32Str()
	actKeyHash, e := bcrypt.GenerateFromPassword([]byte(actKey+actKeySalt), 10)
	if e != nil {
		return &ErrInternal{e}
	}

	// Generate account activation and deactivation URLs.
	// Note that "id.String()" generates a hex string.
	actPart := actKey + "." + id.Hex() + "." + "activate"
	actURL := c.config.DomainName + ActivationEscURL +
		base64.RawURLEncoding.EncodeToString([]byte(actPart))

	deaPart := actKey + "." + id.Hex() + "." + "deactivate"
	deaURL := c.config.DomainName + ActivationEscURL +
		base64.RawURLEncoding.EncodeToString([]byte(deaPart))

	// The activation and deactivation URLs are generated as follows:
	// The key is joined with a 'stringed' chef id, and the word 'activate' or
	// activating, or 'deactivate' for deactivating. These 3 elements are
	// separated with dots.
	// e.g. "veryRandomKeyGoesHere.12345678901234567890.activate"
	// The result is then base64 encoded and appended to the url:
	// "https://recipemanager.io/action/base64EncodedGoodnessGoesHere".

	e = c.SendMail(email,
		"Create Account",
		fmt.Sprintf(NewAccountMessage, actURL, deaURL),
	)
	if e != nil {
		fmt.Println("[ChefsDB.AddChef]", e)
		return &ErrInvalidEmail{email}
	}

	// Add Chef to DB.
	e = c.chefs.Insert(Chef{
		ID:       id,
		Email:    email,
		PwdSalt:  pwdSalt,
		PwdHash:  string(pwdHash),
		Verified: false,
		Admin:    admin,
		Created:  time.Now(),
	})
	if e != nil {
		return &ErrInternal{e}
	}

	// Add Chef activation method.
	e = c.verts.Insert(Verification{
		ID:      id,
		KeySalt: actKeySalt,
		KeyHash: string(actKeyHash),
		Created: time.Now(),
	})
	if e != nil {
		return &ErrInternal{e}
	}
	return nil
}

// SendMail sends mail. Specifically, this is used to send a confirmation email
// to confirm created user.
func (c *ChefsDB) SendMail(to, subject, body string) error {
	to = strings.TrimSpace(to)
	from := c.config.BotEmail
	pass := c.config.BotEmailPwd
	fmt.Println("[ChefsDB.SendMail]", from, pass)

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" + body

	return smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))
}

// MakeActivationEndpoint makes an activation endpoint function to serve using
// http package. The endpoint is used for activating chef accounts and whatnot.
func MakeActivationEndpoint(c *ChefsDB) func(
	http.ResponseWriter, *http.Request) {

	ep := func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.EscapedPath()
		path = path[strings.LastIndex(path, "/")+1:]
		fmt.Println(path)

		page, e := c.ActivateChef(path)
		if e != nil {
			fmt.Println(e)
			http.ServeFile(w, r, "server/static/no-action.html")
			return
		}
		http.ServeFile(w, r, "server/static/"+page+".html")
	}
	return ep
}

// ActivateChef checks escaped url and depending on the validity and data
// contained, it performs appropriate actions.
func (c *ChefsDB) ActivateChef(escapedURLPath string) (page string, e error) {
	decoded, e := base64.RawURLEncoding.DecodeString(escapedURLPath)
	if e != nil {
		return
	}

	// Split escaped path to 3 parts.
	parts := strings.Split(string(decoded), ".")
	if len(parts) != 3 {
		e = &ErrInvalidActivationMethod{
			fmt.Sprintf(
				"got %v parts for escaped path, expected %v",
				len(parts), 3,
			),
		}
		return
	}

	// Obtain activation key, chef id and action.
	key := parts[0]

	if bson.IsObjectIdHex(parts[1]) == false {
		e = &ErrInvalidActivationMethod{
			fmt.Sprintf("unable to parse chefID %v", parts[1]),
		}
		return
	}
	id := bson.ObjectIdHex(parts[1])
	action := parts[2]

	// Check with database.
	v := Verification{}
	e = c.verts.FindId(id).One(&v)
	if e != nil {
		e = &ErrInvalidActivationMethod{e.Error()}
		return
	}

	// Check activation key.
	e = bcrypt.CompareHashAndPassword([]byte(v.KeyHash), []byte(key+v.KeySalt))
	if e != nil {
		e = &ErrInvalidActivationMethod{
			fmt.Sprintf("invalid key: %v", e),
		}
		return
	}

	// Verify/Delete chef based on the specified action.
	switch action {
	case "activate":
		e := c.chefs.UpdateId(id, bson.M{"verified": true})
		if e != nil {
			return "", &ErrInternal{e}
		}
		page = "account-activated"
	case "deactivate":
		e := c.chefs.RemoveId(id)
		if e != nil {
			return "", &ErrInvalidActivationMethod{e.Error()}
		}
		page = "account-deleted"
	}

	// Remove verification entry in db.
	e = c.verts.RemoveId(id)
	return
}

// NewSession (or login) creates a new session for a user. It returns an error
// if not authenticated.
func (c *ChefsDB) NewSession(email, pwd string) (
	sessionInfo *SessionInfo, e error) {

	c.Lock()
	defer c.Unlock()
	email = strings.TrimSpace(email)

	// Find chef with email in database.
	chef := Chef{}
	e = c.chefs.Find(bson.M{"email": email}).One(&chef)
	if e != nil {
		return nil, ErrLoginFailed
	}

	// Check authentication.
	e = bcrypt.CompareHashAndPassword(
		[]byte(chef.PwdHash), []byte(pwd+chef.PwdSalt))
	if e != nil {
		return nil, ErrLoginFailed
	}

	// Generate session id and session key.
	sessionID := bson.NewObjectId()
	key := GetRand32Str()
	keySalt := GetRand32Str()
	keyHash, e := bcrypt.GenerateFromPassword([]byte(key+keySalt), 10)
	if e != nil {
		return nil, &ErrInternal{e}
	}

	// Create Session and SessionInfo objects.
	sessionObj := Session{
		ID:       sessionID,
		ChefID:   chef.ID.Hex(),
		KeySalt:  keySalt,
		KeyHash:  string(keyHash),
		Created:  time.Now(),
		LastSeen: time.Now(),
	}
	sessionInfoObj := SessionInfo{
		SessionID:  sessionID.Hex(),
		SessionKey: key,
		ChefID:     chef.ID.Hex(),
		ChefName:   "",
		ChefEmail:  chef.Email,
	}

	// Add session to chef and update chef in db.
	// TODO: Remove outdated sessions while at it?
	//sessionList := append(chef.Sessions, sessionObj)
	b := bson.M{"$push": bson.M{"sessions": sessionObj}}
	e = c.chefs.UpdateId(chef.ID, b)
	if e != nil {
		return nil, &ErrInternal{e}
	}

	return &sessionInfoObj, nil
}

// DeleteSession Removes a session.
func (c *ChefsDB) DeleteSession(info *SessionInfo) error {
	c.Lock()
	defer c.Unlock()

	if bson.IsObjectIdHex(info.SessionID) == false {
		return &ErrInternal{errors.New("invalid session id - not hex")}
	}
	if bson.IsObjectIdHex(info.ChefID) == false {
		return &ErrInternal{errors.New("invalid chef id - not hex")}
	}

	// Prepare query.
	q := bson.M{
		"$pull": bson.M{
			"sessions": bson.M{
				"_id": bson.ObjectIdHex(info.SessionID)}}}

	// Delete session.
	if e := c.chefs.UpdateId(bson.ObjectIdHex(info.ChefID), q); e != nil {
		return &ErrInternal{e}
	}

	return nil
}

// ClaimSession retrieves a session if authorized.
func (c *ChefsDB) ClaimSession(chefID, sessionID, sessionKey string) (
	sessionInfo *SessionInfo, e error) {

	c.Lock()
	c.Unlock()

	// Check IDs and key.
	if bson.IsObjectIdHex(chefID) == false ||
		bson.IsObjectIdHex(sessionID) == false ||
		sessionKey == "" {
		e = ErrPermissionDenied
		return
	}

	// Obtain session data from db.
	session := Session{}
	q := bson.M{
		"sessions": bson.M{
			"$elemMatch": bson.M{
				"_id": bson.ObjectIdHex(sessionID)}}}
	e = c.chefs.FindId(bson.ObjectIdHex(chefID)).Select(q).One(&session)
	if e != nil {
		e = ErrPermissionDenied
		return
	}

	// Check session_key.
	e = bcrypt.CompareHashAndPassword(
		[]byte(session.KeyHash), []byte(sessionKey+session.KeySalt))
	if e != nil {
		e = ErrPermissionDenied
		return
	}

	// Obtain chef information.
	chef := Chef{}
	e = c.chefs.FindId(bson.ObjectIdHex(chefID)).One(&chef)
	if e != nil {
		e = &ErrInternal{e}
		return
	}

	// Send session info.
	sessionInfo = &SessionInfo{
		SessionID:  session.ID.Hex(),
		SessionKey: sessionKey, // TODO: Send a newly generated sessionKey.
		ChefID:     session.ChefID,
		ChefName:   "",
		ChefEmail:  chef.Email,
	}
	return
}
