package util

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestCheckHappyPath(t *testing.T) {
	assert.Panics(t, func() { Check(errors.New("test")) }, "The code did not panic")
}

func TestCheckWithNilError(t *testing.T) {
	assert.NotPanics(t, func() { Check(nil) })
}

func TestEncryptHappyPath(t *testing.T) {

	secret := "testing"
	encrypted := Encrypt(secret)

	err := bcrypt.CompareHashAndPassword([]byte(encrypted), []byte(secret))
	assert.Nil(t, err)
}

func TestCompareSecretHappyPath(t *testing.T) {

	secret := "testing"
	encrypted, _ := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)

	res := CompareSecret(string(encrypted), secret)
	assert.True(t, res)
}

func TestCompareSecretWithWrongSecret(t *testing.T) {

	secret := "testing"
	encrypted, _ := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)

	res := CompareSecret(string(encrypted), secret+secret)
	assert.False(t, res)
}

func BenchmarkEncrypt(b *testing.B) {
	secret := "testing"
	for i := 0; i < b.N; i++ {
		Encrypt(secret)
	}
}

func TestStrToUintHappyPath(t *testing.T) {
	str := "-100"
	assert.EqualValues(t, -100, StrToUint(str))
}

func BenchmarkStrToUint(b *testing.B) {
	str := "-100"
	for i := 0; i < b.N; i++ {
		StrToUint(str)
	}
}
