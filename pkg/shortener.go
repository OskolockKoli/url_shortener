package shortener

import (
	"crypto/rand"
)

var charset = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_")

const ShortIDLength = 10

type Shortener struct{}

func New() *Shortener {
	return &Shortener{}
}

func (s *Shortener) GenerateShortID() string {
	b := make([]byte, ShortIDLength)
	rand.Read(b)
	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b)
}
