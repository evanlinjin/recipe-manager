package chefs

import (
	"database/sql"
	"sync"

	"encoding/base64"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"strings"
	"net/smtp"
	"net/http"
)

const (
	TableChefs         = "chefs"
	TableVerifications = "verifications"
	TableSessions      = "sessions"
	ActivationEscURL = `/action/`

	NewAccountMessage = `
Hello fellow chef, welcome to Recipe Manager!

To activate your account, click the following link:
%s

If you didn't create an account on Recipe Manager, click here:
%s

Kind Regards,

Team @ Recipe Manager'`
)

type ChefsDB struct {
	sync.Mutex
	db     *sql.DB
	config Config
}

func MakeChefsDB(config Config) (c ChefsDB, e error) {
	c.db, e = sql.Open("sqlite3", "chefs")
	if e != nil {
		return
	}
	return
}

func (c *ChefsDB) Close() error {
	c.Lock()
	defer c.Unlock()
	return c.db.Close()
}

func (c *ChefsDB) Initiate() error {
	c.Lock()
	defer c.Unlock()
	cmd := `CREATE TABLE IF NOT EXISTS %s (
	id       INTEGER PRIMARY KEY,
	email    TEXT NOT NULL UNIQUE,
	pwd_salt TEXT NOT NULL,
	pwd_hash TEXT NOT NULL,
	verified BOOL DEFAULT 0 NOT NULL,
	created  TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL)`
	if _, e := c.db.Exec(fmt.Sprintf(cmd, TableChefs)); e != nil {
		return e
	}
	cmd = `CREATE TABLE IF NOT EXISTS %s (
	id       INTEGER PRIMARY KEY,
	key_salt TEXT NOT NULL,
	key_hash TEXT NOT NULL,
	created  TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL)`
	if _, e := c.db.Exec(fmt.Sprintf(cmd, TableVerifications)); e != nil {
		return e
	}
	cmd = `CREATE TABLE IF NOT EXISTS %s (
	id        INTEGER PRIMARY KEY,
	chef_id   INTEGER NOT NULL,
	key_salt  TEXT NOT NULL,
	key_hash  TEXT NOT NULL,
	created   TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	last_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL)`
	if _, e := c.db.Exec(fmt.Sprintf(cmd, TableSessions)); e != nil {
		return e
	}
	return nil
}

func (c *ChefsDB) AddChef(email, pwd string) error {
	c.Lock()
	defer c.Unlock()
	email = strings.TrimSpace(email)

	// Check if chef already exists with specified email.
	scannedEmail := ""
	c.db.QueryRow("SELECT email FROM chefs WHERE email = ?", email).
		Scan(&scannedEmail)
	if scannedEmail == email {
		return &ErrChefAlreadyExists{email}
	}

	// Generate chef id, pwd_salt and pwd_hash.
	id := GetRandUniqID()
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
	actPart := actKey + "." + strconv.FormatInt(id, 10) + "." + "activate"
	actURL := c.config.DomainName + ActivationEscURL +
		base64.RawURLEncoding.EncodeToString([]byte(actPart))

	deaPart := actKey + "." + strconv.FormatInt(id, 10) + "." + "deactivate"
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
		return &ErrInvalidEmail{email}
	}

	// Add Chef to DB.
	_, e = c.db.Exec(
		"INSERT INTO chefs(id,email,pwd_salt,pwd_hash) VALUES(?,?,?,?)",
		id, email, pwdSalt, pwdHash,
	)
	if e != nil {
		return &ErrInternal{e}
	}

	// Add Chef activation method.
	_, e = c.db.Exec(
		"INSERT OR REPLACE INTO verifications(id,key_salt,key_hash) VALUES(?,?,?)",
		id, actKeySalt, actKeyHash,
	)
	if e != nil {
		return &ErrInternal{e}
	}
	return nil
}

// SendMail sends mail.
func (c *ChefsDB) SendMail(to, subject, body string) error {
	to = strings.TrimSpace(to)
	from := c.config.BotEmail
	pass := c.config.BotEmailPwd

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

		if e := c.ActivateChef(path); e != nil {
			fmt.Println(e)
			http.ServeFile(w, r, "server/static/account-activation-fail.html")
			return
		}
		http.ServeFile(w, r, "server/static/account-activation-success.html")
	}
	return ep
}



// ActivateChef checks escaped url and depending on the validity and data
// contained, it performs appropriate actions.
func (c *ChefsDB) ActivateChef(escapedURLPath string) error {
	decoded, e := base64.RawURLEncoding.DecodeString(escapedURLPath)
	if e != nil {
		return e
	}

	// Split escaped path to 3 parts.
	parts := strings.Split(string(decoded), ".")
	if len(parts) != 3 {
		return &ErrInvalidActivationMethod{
			fmt.Sprintf(
				"got %v parts for escaped path, expected %v",
				len(parts), 3,
			),
		}
	}

	// Obtain activation key, chef id and action.
	key := parts[0]
	id, e := strconv.ParseInt(parts[1], 10, 64)
	if e != nil {
		return &ErrInvalidActivationMethod{
			fmt.Sprintf(
				"unable to parse chefID %v because %v",
				parts[1], e,
			),
		}
	}
	action := parts[2]

	// Check with database.
	keySalt, keyHash := "", ""
	e = c.db.QueryRow(
		"SELECT key_salt,key_hash FROM verifications WHERE id = ?", id,
	).Scan(&keySalt, &keyHash)
	if e != nil {
		return &ErrInvalidActivationMethod{e.Error()}
	}

	// Check activation key.
	e = bcrypt.CompareHashAndPassword([]byte(keyHash), []byte(key+keySalt))
	if e != nil {
		return &ErrInvalidActivationMethod{
			fmt.Sprintf("invalid key: %v", e),
		}
	}

	// Verify/Delete chef based on the specified action.
	switch action {
	case "activate":
		_, e := c.db.Exec(
			"UPDATE chefs SET verified = ? WHERE id = ?", true, id,
		)
		if e != nil {
			return &ErrInternal{e}
		}
	case "deactivate":
		_, e := c.db.Exec(
			"DELETE FROM chefs WHERE id = ? AND verified = ?", id, false,
		)
		if e != nil {
			return &ErrInvalidActivationMethod{e.Error()}
		}
	}

	// Remove verification entry in db.
	_, e = c.db.Exec("DELETE FROM verifications WHERE id = ?", id)
	return e
}

