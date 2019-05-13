package main

import (
	"auth/dynamo"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

const secondsForOneHour = 3600
const hmacSampleSecret = "GYthbtJJ6tp3852JMEVmVHhDckdHHDsJ"

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

func getJWT(user dynamo.User, relation string) string {

	now := time.Now().Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_name": user.UserName,
		"role_name": user.RoleName,
		"relation":  relation,
		"exp":       now + secondsForOneHour,
		"nbf":       now,
		"iat":       now,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(hmacSampleSecret))
	check(err)

	return tokenString
}

func verifyJWT(tokenString string) (err error, roleName, relation string) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(hmacSampleSecret), nil
	})

	if err != nil {
		return err, "", ""
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		roleName, roleNameOK := claims["role_name"].(string)
		if !roleNameOK {
			return errors.New("Fail to convert interface to string"), "", ""
		}
		relation, relationOK := claims["relation"].(string)
		if !relationOK {
			return errors.New("Fail to convert interface to string"), "", ""
		}

		return nil, string(roleName), string(relation)
	}
	return errors.New("Unexpected invalid token"), "", ""
}
