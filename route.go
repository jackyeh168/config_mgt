package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const hmacSampleSecret = "GYthbtJJ6tp3852JMEVmVHhDckdHHDsJ"

func saveUser(user UserInfo) bool {
	db := getDBInstance()

	if db.Where("username = ?", user.Username).Find(&user).RecordNotFound() {
		user.Password = encrypt(user.Password)
		check(db.Create(&user).Error)
		return true
	}
	return false
}

func verifyUser(user UserInfo) (bool, UserInfo) {
	db := getDBInstance()
	inputPassword := user.Password

	if db.Where(&UserInfo{Username: user.Username}).First(&user).RecordNotFound() {
		return false, UserInfo{}
	} else {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(inputPassword)); err != nil {
			return false, UserInfo{}
		}
		return true, user
	}
}

func register() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := UserInfo{}
		err := c.BindJSON(&user)
		check(err)

		if saveUser(user) {
			c.JSON(200, gin.H{})
		} else {
			c.Status(http.StatusBadRequest)
		}
		return
	}
}

func login() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := UserInfo{}
		err := c.BindJSON(&user)
		check(err)
		if isValid, user := verifyUser(user); isValid {
			c.JSON(200, gin.H{
				"token": getJWT(user),
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

		if roleID, err := verifyJWT(token); err == nil {
			// if success, then next
			c.Set("roleID", roleID)
			c.Next()
		} else {
			// else return http 401
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		return
	}
}

func getRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/login", login())
	r.POST("/register", register())

	authorized := r.Group("/")
	// per group middleware! in this case we use the custom created
	// AuthRequired() middleware just in the "authorized" group.
	authorized.Use(authRequired())
	{
		authorized.GET("/getData", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"data": "justfortest",
			})
		})

		// authorized.GET("/users")
		// authorized.POST("/user") // add user

		// authorized.GET("/user/:user_id/projects")
		// authorized.POST("/user/:user_id/project") // assign project for user

		// authorized.GET("/project/:project_id/envs")
		// authorized.POST("/project/:project_id/env") // add env for project

		// authorized.GET("/projects")
		// authorized.POST("/project") // add project
	}
	return r
}
