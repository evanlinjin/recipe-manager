package handle

import (
	"encoding/json"
	"github.com/evanlinjin/recipe-manager/server/chefs"
)

type LoginResponse struct {
	Okay    bool               `json:"okay"`
	Session *chefs.SessionInfo `json:"session,omitempty"`
	Message string             `json:"message,omitempty"`
}

type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func GetLoginData(in interface{}) (out LoginData, e error) {
	data, e := json.Marshal(in)
	if e != nil {
		return
	}
	e = json.Unmarshal(data, &out)
	return
}

func HandleNewSessionError(e error) string {
	switch e.(type) {
	case *chefs.ErrInternal:
		return "Server error while creating new session."
	default:
		return "Invalid request."
	}
}
