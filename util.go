package main

import "math/rand"

func generateSlug() string {
	urlChars := []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	s := make([]rune, 4)
	for i := range s {
		s[i] = urlChars[rand.Intn(len(urlChars))]
	}
	return string(s)
}
