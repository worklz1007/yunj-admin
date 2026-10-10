package member

import (
	"crypto/rand"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 10

const passwordAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func HashPassword(raw string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(raw), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func VerifyPassword(hash, raw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(raw)) == nil
}

func RandomPassword() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	out := make([]byte, len(buf))
	for i, b := range buf {
		out[i] = passwordAlphabet[int(b)%len(passwordAlphabet)]
	}
	if utf8.RuneCount(out) != 16 {
		return "", ErrInternal
	}
	return string(out), nil
}
