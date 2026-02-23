package utils

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func GeneratePassword() string {
	rand.Seed(time.Now().UnixNano())
	chars := "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz23456789@#"
	pass := make([]byte, 10)
	for i := range pass {
		pass[i] = chars[rand.Intn(len(chars))]
	}
	return string(pass)
}

func HashPassword(p string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(p), 10)
	return string(b), err
}
