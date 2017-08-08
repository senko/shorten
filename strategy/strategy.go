package strategy

import (
	"../store"
	"math/rand"
)

type GenerateKey func(store.Store) string

const (
	DefaultAlphabet = "abcdefghijklmnopqrstuvwxyz0123456789"
	DefaultLength   = 6
)

func RandomKey(alphabet string, length int) GenerateKey {
	if alphabet == "" {
		alphabet = DefaultAlphabet
	}
	runes := []rune(alphabet)

	if length == 0 {
		length = DefaultLength
	}

	return func(s store.Store) string {
		buf := make([]rune, length)
		for i := range buf {
			buf[i] = runes[rand.Intn(len(runes))]
		}
		return string(buf)
	}
}

func DefaultRandomKey() GenerateKey {
	return RandomKey(DefaultAlphabet, DefaultLength)
}
