package chefs

import "time"

// Config contains data to configure ChefsDB.
type Config struct {
	DomainName string
	BotEmail string
	BotEmailPwd string
}

type Chef struct {
	ID       int64
	Email    string
	PwdSalt  string
	PwdHash  string
	Verified bool
	Created  time.Time
}

type Verification struct {
	ID      int64
	KeySalt string
	KeyHash string
	Created time.Time
}

type Session struct {
	ID       int64
	ChefID   string
	KeySalt  string
	KeyHash  string
	Meta     string
	Created  time.Time
	LastSeen time.Time
}
