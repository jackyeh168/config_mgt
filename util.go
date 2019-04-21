package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

const secondsForOneHour = 3600

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func encrypt(secret string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	check(err)
	return string(hash)
}

func getJWT(user UserInfo) string {

	now := time.Now().Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"role_id":  user.RoleID,
		"exp":      now + secondsForOneHour,
		"nbf":      now,
		"iat":      now,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(hmacSampleSecret))
	check(err)

	return tokenString
}

func verifyJWT(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return 0, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(hmacSampleSecret), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		v, ok := claims["role_id"].(float64)
		if !ok {
			return 0, errors.New("Fail to convert interface to uint")
		}
		return uint(v), nil
	}
	return 0, errors.New("Unexpected invalid token")
}
