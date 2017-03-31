package chefs

import (
	"net/http"
	"strings"
)

func MakeActivationEndpoint(c *ChefsDB) func(
	http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.EscapedPath()
		path = path[strings.LastIndex(path, "/")+1:]

	}
}