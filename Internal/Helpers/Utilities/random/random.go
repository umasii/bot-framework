package utils

import (
	"math/rand"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int, customAlphabet string) string {
	runeSet := letterRunes
	if customAlphabet != "" {
		runeSet = []rune(customAlphabet)
	}

	b := make([]rune, n)
	rand.Seed(time.Now().UnixNano())

	for i := range b {
		b[i] = runeSet[rand.Intn(len(runeSet))]
	}
	return string(b)
}
