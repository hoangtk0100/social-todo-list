package common

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSequence(n int) string {
	b := make([]rune, n)

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	for index := range b {
		b[index] = letters[r.Intn(999999)%len(letters)]
	}

	return string(b)
}

func GenSalt(length int) string {
	if length < 0 {
		length = 50
	}

	return randSequence(length)
}

type md5Hash struct{}

func NewMd5Hash() *md5Hash {
	return &md5Hash{}
}

func (h *md5Hash) Hash(data string) string {
	hasher := md5.New()
	hasher.Write([]byte(data))
	return hex.EncodeToString(hasher.Sum(nil))
}
