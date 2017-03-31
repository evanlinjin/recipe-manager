package chefs

import (
	"net/http"
	//"strings"
	"fmt"
	"strings"
)

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
