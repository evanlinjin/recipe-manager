package talkrelay

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// MakePackage makes a package to send via websocket for client/server talk.
func MakePackage(obj interface{}, key []byte) (pac []byte, e error) {
	key = bytes.TrimSpace(key)

	objBytes, e := json.Marshal(obj)
	if e != nil {
		return
	}

	var dataPart = base64.RawURLEncoding.EncodeToString(objBytes)

	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(dataPart))

	var sigPart = base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	pac = []byte(dataPart + "." + sigPart)
	return
}

// ReadPackage reads a package.
func ReadPackage(pac []byte, key []byte, obj interface{}) error {
	fmt.Println("[ReadPackage] Pkg:", string(pac), ", Key:", string(key))
	dot := bytes.Index(pac, []byte("."))
	if dot == -1 {
		return fmt.Errorf("invalid package - dot at %v", dot)
	}

	var dataPart = pac[:dot]
	var sigPart = pac[dot+1:]

	mac := hmac.New(sha256.New, key)
	mac.Write(dataPart)

	var genSig = mac.Sum(nil)
	obtSig, e := base64.RawURLEncoding.DecodeString(string(sigPart))
	if e != nil {
		return e
	}
	if hmac.Equal(genSig, []byte(obtSig)) == false {
		return fmt.Errorf("unmatched '%s' and '%s'", genSig, obtSig)
	}
	decData, e := base64.RawURLEncoding.DecodeString(string(dataPart))
	if e != nil {
		return e
	}

	json.Unmarshal(decData, obj)
	return nil
}
