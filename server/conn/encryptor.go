package conn

import (
	"encoding/base64"

	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"sync"
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

// Encrypt encrypts plainText to cipherText, then base64-encodes it to become
// an encodedCipher. Note that the iv is prepended on cipherText - this is
// collectively referred to as cipherText.
func (r *Encryptor) Encrypt(plainText []byte) (encodedCipher []byte, e error) {
	for len(plainText)%DefSize != 0 {
		plainText = append(plainText, byte(' '))
	}

	r.Lock()
	block, e := aes.NewCipher(r.Key)
	r.Unlock()
	if e != nil {
		return
	}

	// Consists of [iv|cipherText].
	cipherText := make([]byte, DefSize+len(plainText))

	// Make iv part.
	iv := cipherText[:DefSize]
	if _, e := io.ReadFull(rand.Reader, iv); e != nil {
		return nil, e
	}

	cipher.NewCBCEncrypter(block, iv).CryptBlocks(cipherText[DefSize:], plainText)
	encodedCipher = []byte(base64.RawURLEncoding.EncodeToString(cipherText))
	return
}

// Decrypt decrypts an encodedCipher to plainText. It assumes that the
// cipherText derived from the encodedCipher has it's iv prepended to it.
func (r *Encryptor) Decrypt(encodedCipher []byte) (plainText []byte, e error) {
	cipherText, e := base64.RawURLEncoding.DecodeString(string(encodedCipher))

	r.Lock()
	block, e := aes.NewCipher(r.Key)
	r.Unlock()
	if e != nil {
		return
	}

	// Split into iv|cipherText parts.
	if len(cipherText) < DefSize {
		e = fmt.Errorf("cipherText too short at %v, expecting >= %v",
			len(cipherText), DefSize)
		return
	}
	iv, cipherText := cipherText[:DefSize], cipherText[DefSize:]

	// Check if satisfies block size.
	if len(cipherText)%DefSize != 0 {
		e = fmt.Errorf("cipherText should be multiple of %v, got %v",
			DefSize, len(cipherText))
		return
	}

	plainText = make([]byte, len(cipherText))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(plainText, cipherText)
	plainText = bytes.TrimSpace(plainText)
	return
}
