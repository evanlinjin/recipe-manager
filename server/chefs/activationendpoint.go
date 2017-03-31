package chefs

import "net/http"

func MakeActivationEndpoint(chefsDB *ChefsDB) func(
	http.ResponseWriter, *http.Response) {
	return func(w http.ResponseWriter, r *http.Response) {

	}
}