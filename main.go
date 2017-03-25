package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
		w.WriteHeader(http.StatusOK)
		return
	})
	// http.ListenAndServe(GetPort(), nil)
	http.ListenAndServe(":8080", nil)
}

// GetPort returns the port number.
func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	return ":" + port
}
