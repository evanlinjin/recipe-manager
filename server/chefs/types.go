package chefs

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	DBAuth  = "auth"
	CTChefs = "chefs"
	CTVerts = "verifications"

	ActivationEscURL  = `/action/`
	NewAccountMessage = `
Hello fellow chef, welcome to Recipe Manager!

To activate your account, click the following link:
%s

If you didn't create an account on Recipe Manager, click here:
%s

Kind Regards,

Team @ Recipe Manager'`
)

// Config contains data to configure ChefsDB.
type Config struct {
	DomainName  string
	MongoUrls   string // separated by commas.
	BotEmail    string
	BotEmailPwd string
}

// Chef represents how a Chef is recorded in database.
type Chef struct {
	ID       bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Email    string        `json:"email"`
	PwdSalt  string        `json:"pwd_salt"`
	PwdHash  string        `json:"pwd_hash"`
	Verified bool          `json:"vertified"`
	Created  time.Time     `json:"created"`
	Config   ChefConfig    `json:"config,omitempty"`
	Sessions []Session     `json:"sessions"`
	Teams    []string      `json:"teams,omitempty"`
}

// Verification represents how a verification method is recorded in database.
type Verification struct {
	ID      bson.ObjectId `json:"id" bson:"_id,omitempty"`
	KeySalt string        `json:"key_salt"`
	KeyHash string        `json:"key_hash"`
	Created time.Time     `json:"created"`
}

// Session represents how a session is recorded in database.
type Session struct {
	ID       bson.ObjectId `json:"id" bson:"_id,omitempty"`
	ChefID   string        `json:"chef_id"`
	KeySalt  string        `json:"key_salt"`
	KeyHash  string        `json:"key_hash"`
	Meta     string        `json:"meta"`
	Created  time.Time     `json:"created"`
	LastSeen time.Time     `json:"last_seen"`
}

// SessionInfo represents how a session is to be presented to the client.
type SessionInfo struct {
	SessionID  string   `json:"session_id"`
	SessionKey string   `json:"session_key"`
	ChefID     string   `json:"user_id"`
	ChefName   string   `json:"chef_name"`
	ChefEmail  string   `json:"chef_email"`
	Teams      []string `json:"teams"`
}

// ChefConfig is yet to be implemented.
type ChefConfig struct {
}
