package random

import (
	"math/rand"
	"time"
)

var initialised = false
var validchars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GenerateRandomString(n int) string {
	if !initialised {
		rand.Seed(time.Now().UnixNano())
		initialised = true
	}
	b := make([]rune, n)
	for i := range b {
		b[i] = validchars[rand.Intn(len(validchars))]
	}
	return string(b)
}
