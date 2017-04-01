package handle

import "encoding/json"

type NewChefData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func GetNewChefData(in interface{}) (out NewChefData, e error) {
	data, e := json.Marshal(in)
	if e != nil {
		return
	}
	e = json.Unmarshal(data, &out)
	return
}