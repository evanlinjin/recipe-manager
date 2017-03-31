package chefs

import (
	"database/sql"
	"sync"

	"encoding/base64"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

const (
	TableChefs         = "chefs"
	TableVerifications = "verifications"
	TableSessions      = "sessions"
)

type ChefsDB struct {
	sync.Mutex
	db *sql.DB
}

func MakeChefsDB() (c ChefsDB, e error) {
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

	// Check if chef already exists with specified email.
	verified := false
	c.db.QueryRow(
		"SELECT verified FROM chefs WHERE email = ?").Scan(&verified)
	if verified == true {
		return &ErrChefAlreadyExists{email}
	}

	// Generate id, pwd_salt and pwd_hash.
	id := GetRandUniqID()
	pwdSalt := GetRand32Str()
	pwdHash, e := bcrypt.GenerateFromPassword([]byte(pwd+pwdSalt), 10)
	if e != nil {
		return &ErrInternal{e}
	}

	// TODO: Send confirmation email, and check email validity.
	e = SendMail(email, "Create Account", "Testing yoyoyoy!")
	if e != nil {
		return e
	}

	// Add Chef to DB.
	_, e = c.db.Exec(
		"INSERT INTO chefs(id,email,pwd_salt,pwd_hash) VALUES (?,?,?,?)",
		id, email, pwdSalt, pwdHash,
	)
	return e
}

func GetRandUniqID() int64 {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Int63()
}

func GetRand32Str() string {
	b := make([]byte, ((3*32)/4)+1)
	rand.New(rand.NewSource(time.Now().UnixNano())).Read(b)
	return base64.RawURLEncoding.EncodeToString(b)[:32]
}
