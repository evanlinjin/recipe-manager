package conn

import (
	"encoding/base64"

	"crypto/rand"
	"fmt"
	"io"
	"sync"
	"crypto/aes"
)

const DefSize = aes.BlockSize

// Encryptor is responsible for encryption and decryption. Specifically, it is
// used for encrypting and decrypting packages and it's signature. The key
// member is a client/server agreed key that is used to encrypt and decrypt.
type Encryptor struct {
	sync.Mutex
	Key []byte
}

// MakeEncryptor makes an Encryptor instance.
func MakeEncryptor() (enc Encryptor) {
	enc.Key = make([]byte, DefSize)
	for i, _ := range enc.Key {
		enc.Key[i] = byte(0)
	}
	return
}

// Makes a random key (does not set it).
func (r *Encryptor) makeKey() ([]byte, error) {
	key := make([]byte, DefSize)
	if _, e := io.ReadFull(rand.Reader, key); e != nil {
		return nil, e
	}
	return []byte(base64.RawURLEncoding.EncodeToString(key)), nil
}

// Sets a key.
func (r *Encryptor) setKey(encKey []byte) (e error) {
	r.Lock()
	defer r.Unlock()
	key, e := base64.RawURLEncoding.DecodeString(string(encKey))
	if e != nil {
		return
	}
	if len(key) != DefSize {
		return fmt.Errorf("length of key should be %v, got %v", DefSize, len(key))
	}
	r.Key = key
	return
}

// Encrypt really only encodes.
func (r *Encryptor) Encrypt(plainText []byte) (encodedCipher []byte, e error) {
	encodedCipher = []byte(base64.RawURLEncoding.EncodeToString(plainText))
	return
}

// Decrypt really only decodes.
func (r *Encryptor) Decrypt(encodedCipher []byte) (plainText []byte, e error) {
	plainText, e = base64.RawURLEncoding.DecodeString(string(encodedCipher))
	return
}
