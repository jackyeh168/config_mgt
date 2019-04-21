package main

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const hmacSampleSecret = "GYthbtJJ6tp3852JMEVmVHhDckdHHDsJ"

func verifyUser(user UserInfo) bool {
	return true
}

func generateJWT() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"foo": "bar",
		"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(hmacSampleSecret)
	check(err)

	return tokenString
}

func verifyJWT(jwt string) bool {
	return true
}

func login() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := UserInfo{}
		err := c.BindJSON(&user)
		check(err)

		if verifyUser(user) {
			c.JSON(200, gin.H{
				"token": generateJWT(),
			})
		} else {
			c.Status(http.StatusUnauthorized)
		}
		return
	}
}

func authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.GetHeader("Authorization")
		// Verify jwt token
		if verifyJWT(token) {
			// if success, then next
			c.Next()
		} else {
			// else return http 401
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		return
	}
}
