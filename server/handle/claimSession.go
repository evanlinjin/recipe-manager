package handle

import (
	"encoding/json"
	"github.com/evanlinjin/recipe-manager/server/chefs"
)

type ClaimSessionData struct {
	ChefID     string `json:"chef_id"`
	SessionID  string `json:"session_id"`
	SessionKey string `json:"session_key"`
}

func GetClaimSessionData(in interface{}) (out ClaimSessionData, e error) {
	data, e := json.Marshal(in)
	if e != nil {
		return
	}
	e = json.Unmarshal(data, &out)
	return
}

func DealClaimSessionError(e error) string {
	switch e.(type) {
	case *chefs.ErrInternal:
		return "Server error while claiming session."
	default:
		return "Your session has expired."
	}
}