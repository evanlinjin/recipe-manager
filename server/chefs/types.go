package chefs

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	DBAuth  = "auth"
	CTChefs = "chefs"
	CTVerts = "verifications"
)

// Config contains data to configure ChefsDB.
type Config struct {
	DomainName  string
	MongoUrls   string // separated by commas.
	BotEmail    string
	BotEmailPwd string
}

type Chef struct {
	ID       bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Email    string        `json:"email"`
	PwdSalt  string        `json:"pwd_salt"`
	PwdHash  string        `json:"pwd_hash"`
	Verified bool          `json:"vertified"`
	Created  time.Time     `json:"created"`
}

type Verification struct {
	ID      bson.ObjectId `json:"id" bson:"_id,omitempty"`
	KeySalt string        `json:"key_salt"`
	KeyHash string        `json:"key_hash"`
	Created time.Time     `json:"created"`
}

type Session struct {
	ID       bson.ObjectId `json:"id" bson:"_id,omitempty"`
	ChefID   string
	KeySalt  string
	KeyHash  string
	Meta     string
	Created  time.Time
	LastSeen time.Time
}
