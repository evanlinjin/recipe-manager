package talkrelay

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// Packages, in Recipe Manager terms, is a encoded and signed message. This is
// to, hopefully, ensure authenticity of messages transported between server and
// client.
//
// TODO: Describe the package structure.
//
// Note that packages are to be encrypted with the Encryptor, and appended to an
// encrypted signature key (used to sign the package). These two parts are
// separated with a dot. The result is an encryptedPackage (or cipherPackage)
// ready to be sent between server and client.

// MakePackage makes a package with a given signature key.
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

// ReadPackage reads a decoded package and checks it's authenticity with it's
// decoded signature key. ReadPackage assumes that the message is in a json
// format, and unmarshals the result to obj if package is verified.
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
