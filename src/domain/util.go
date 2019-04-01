package domain

import (
	"math/rand"

	"golang.org/x/crypto/bcrypt"
)

// Encrypt encrypts the given password
func Encrypt(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hash)
}

// GenerateSlug generates slug
func GenerateSlug(length int) string {
	urlChars := []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	s := make([]rune, length)
	for i := range s {
		s[i] = urlChars[rand.Intn(len(urlChars))]
	}
	return string(s)
}
