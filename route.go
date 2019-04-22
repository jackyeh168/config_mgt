package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const hmacSampleSecret = "GYthbtJJ6tp3852JMEVmVHhDckdHHDsJ"

type jsonPermission struct {
	User    []string `json: "user"`
	Project []string `json: "project"`
	Env     []string `json: "env"`
}

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

func contains(arr []string, element string) bool {
	for _, v := range arr {
		if element == v {
			return true
		}
	}
	return false
}

func verifyPermission(roleID, resource, action string) bool {
	db := getDBInstance()

	role := RoleInfo{}
	check(db.Where("id =?", roleID).Find(&role).Error)
	var j jsonPermission

	err := json.Unmarshal([]byte(role.Permission), &j)
	check(err)

	v := reflect.ValueOf(&j).Elem().FieldByName(strings.Title(resource))
	arr := v.Interface().([]string)
	// fmt.Println(arr)

	return contains(arr, action)
}

func isValid(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verify role
		roleID, ok := c.Get("roleID")
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		roleStr := fmt.Sprint(roleID.(uint))
		if verifyPermission(roleStr, resource, action) {
			// if permission is valid, then next
			c.Next()
			return
		}
		// else return http 401
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}

func getUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := getDBInstance()
		var users []UserInfo
		db.Select("id, username, role_id").Find(&users)

		c.JSON(200, gin.H{
			"data": users,
		})
	}
}

func addUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		// db := getDBInstance()
		var users []UserInfo
		c.BindJSON(&users)

		c.JSON(200, gin.H{
			"data": users,
		})
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

		authorized.GET("/users", isValid("user", "read"), getUsers())
		authorized.POST("/users", isValid("user", "create"), addUsers()) // add user
		// authorized.POST("/user/:user_id/project", isValid("project", "create"), addProjectsForUser()) // assign project for user

		// authLeader := authorized.Group("leader")
		// authLeader.Use(isLeader())
		// {
		// 	authLeader.GET("/user/:user_id/projects")
		// 	authLeader.POST("/project/:project_id/env") // add env for project
		// }

		// authorized.POST("/user") // add user

		// authorized.GET("/user/:user_id/projects")

		// authorized.GET("/project/:project_id/envs")
		// authorized.POST("/project/:project_id/env") // add env for project

		// authorized.GET("/projects")
		// authorized.POST("/project") // add project
	}
	return r
}
