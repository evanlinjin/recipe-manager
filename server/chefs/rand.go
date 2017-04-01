package chefs

import (
	"time"
	"math/rand"
	"encoding/base64"
)

func GetRandUniqID() int64 {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Int63()
}

func GetRand32Str() string {
	b := make([]byte, ((3*32)/4)+1)
	rand.New(rand.NewSource(time.Now().UnixNano())).Read(b)
	return base64.RawURLEncoding.EncodeToString(b)[:32]
}