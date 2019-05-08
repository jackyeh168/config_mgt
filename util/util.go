package util

import (
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

// Check is to check error
func Check(err error) {
	if err != nil {
		panic(err)
	}
}

// Encrypt is a function use bycrpyt
func Encrypt(secret string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	Check(err)
	return string(hash)
}

// CompareSecret is a function to compare the encrypted string and plain text
func CompareSecret(encrypted, secret string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encrypted), []byte(secret))
	if err != nil {
		return false
	}
	return true
}

// StrToUint is the function that convert String to uint64
func StrToUint(s string) uint {
	i, err := strconv.Atoi(s)
	Check(err)
	return (uint)(i)
}
