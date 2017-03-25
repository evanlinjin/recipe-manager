package main

import (
	"log"
	"net/http"
	"os"
	"time"

	mgo "gopkg.in/mgo.v2"

	"github.com/kabukky/httpscerts"
)

const (
	_StaticIP = "34.204.161.180"
	_Port     = "8080"
)

func main() {
	if e := httpscerts.Check("cert.pem", "key.pem"); e != nil {
		if e = httpscerts.Generate("cert.pem", "key.pem", _StaticIP+":"+_Port); e != nil {
			log.Fatal("Error: Couldn't create https certs.")
		}
	}

	http.HandleFunc("/", handler)
	http.ListenAndServeTLS(":"+_Port, "cert.pem", "key.pem", nil)
}

// Timestamp records a timestamp
type Timestamp struct {
	epoch int64
}

// NewTimestamp creates a new timestamp.
func NewTimestamp() *Timestamp {
	return &Timestamp{
		epoch: time.Now().UnixNano(),
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	session, e := mgo.Dial("mongodb://127.0.0.1:32017")
	if e != nil {
		panic(e)
	}
	defer session.Close()

	c := session.DB("test").C("timestaps")
	if e = c.Insert(NewTimestamp()); e != nil {
		log.Fatal(e)
	}

	n, e := c.Count()
	if e != nil {
		log.Fatal(e)
	}

	w.Write([]byte("Hello World! We have had " + string(n) + " visits!"))
	w.WriteHeader(http.StatusOK)
	return
}

// GetPort returns the port number.
func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	return ":" + port
}
