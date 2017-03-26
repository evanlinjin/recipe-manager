package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/evanlinjin/recipe-manager/config"
)

const (
	_StaticIP   = "34.204.161.180"
	_DomainName = "recipemanager.io"
	_Port       = "8080"
	_PortHTTP   = "80"
)

func main() {
	// Get configuration.
	c, e := config.GetNetworkConfig()
	if e != nil {
		log.Fatal(e)
	}

	http.HandleFunc("/", MakeHander(c))

	// Listen and serve.
	log.Printf("Listening and serving on %s:%s ...", c.Domain, c.Port)
	e = http.ListenAndServeTLS(":"+c.Port, c.SSLCertPath, c.SSLKeyPath, nil)
	if e != nil {
		log.Fatal(e)
	}
}

// MakeHandler makes an http handler.
func MakeHander(c *config.NetworkConfig) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Redirect HTTP.
		if r.TLS == nil {
			http.Redirect(w, r, "https://"+c.Domain, http.StatusMovedPermanently)
			return
		}

		w.Write([]byte("Hello World!"))
		w.WriteHeader(http.StatusOK)
		return
	}
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
	w.Write([]byte("Hello World!"))
	w.WriteHeader(http.StatusOK)
	return
	//if r.TLS == nil {
	//	http.Redirect(w, r, "https://"+_DomainName, http.StatusMovedPermanently)
	//	return
	//}
	//
	//w.Header().Set("Content-Type", "text/html;charset=UTF-8")
	//session, e := mgo.Dial("mongodb://127.0.0.1:32017")
	//if e != nil {
	//	w.Write([]byte(e.Error()))
	//	w.WriteHeader(http.StatusInternalServerError)
	//	return
	//}
	//defer session.Close()
	//
	//c := session.DB("test").C("timestaps")
	//if e = c.Insert(NewTimestamp()); e != nil {
	//	w.Write([]byte(e.Error()))
	//	w.WriteHeader(http.StatusInternalServerError)
	//	return
	//}
	//
	//n, e := c.Count()
	//if e != nil {
	//	w.Write([]byte(e.Error()))
	//	w.WriteHeader(http.StatusInternalServerError)
	//	return
	//}
	//
	//w.Write([]byte("Hello World! We have had " + strconv.Itoa(n) + " visits!"))
	//w.WriteHeader(http.StatusOK)
	//return
}

// GetPort returns the port number.
func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	return ":" + port
}
