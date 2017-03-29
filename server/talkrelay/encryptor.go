package talkrelay

import (
	"encoding/base64"

	"fmt"
	"crypto/aes"
	"io"
	"crypto/rand"
	"crypto/cipher"
)

const DefSize = aes.BlockSize

type Encryptor struct {
	Key []byte
}

func MakeEncryptor() (enc Encryptor) {
	enc.Key = make([]byte, DefSize)
	for i, _ := range enc.Key {
		enc.Key[i] = byte(0)
	}
	return
}

func (r *Encryptor) MakeRandomKey() ([]byte, error) {
	r.Key = make([]byte, DefSize)
	if _, e := io.ReadFull(rand.Reader, r.Key); e != nil {
		return nil, e
	}
	return []byte(base64.RawURLEncoding.EncodeToString(r.Key)), nil
}

func (r *Encryptor) ReadEncodedKey(encKey []byte) (e error) {
	key , e := base64.RawURLEncoding.DecodeString(string(encKey))
	if e != nil {
		return
	}
	if len(key) != DefSize {
		return fmt.Errorf("length of key should be %v, got %v", DefSize, len(key))
	}
	r.Key = key
	return
}

func (r *Encryptor) Encrypt(plainText []byte) (encodedCipher []byte, e error) {
	for len(plainText)%DefSize != 0 {
		plainText = append(plainText, byte(' '))
	}

	block, e := aes.NewCipher(r.Key)
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

func (r *Encryptor) Decrypt(encodedCipher []byte) (plainText []byte, e error) {
	cipherText, e := base64.RawURLEncoding.DecodeString(string(encodedCipher))

	block, e := aes.NewCipher(r.Key)
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
	if len(cipherText) % DefSize != 0 {
		e = fmt.Errorf("cipherText should be multiple of %v, got %v",
			DefSize, len(cipherText))
		return
	}

	plainText = make([]byte, len(cipherText))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(plainText, cipherText)
	return
}