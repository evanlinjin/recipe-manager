package chefs

import "time"

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
