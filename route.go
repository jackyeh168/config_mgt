package main

import (
	"auth/controller"
	"auth/model"
	"auth/util"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type jsonPermission struct {
	User    []string `json: "user"`
	Project []string `json: "project"`
	Env     []string `json: "env"`
}

func verifyUser(user model.UserInfo) (bool, model.UserInfo) {
	db := model.GetDBInstance()
	inputPassword := user.Password

	if db.Where(&model.UserInfo{Username: user.Username}).First(&user).RecordNotFound() {
		return false, model.UserInfo{}
	} else {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(inputPassword)); err != nil {
			return false, model.UserInfo{}
		}
		return true, user
	}
}

// func register() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		user := model.UserInfo{}
// 		err := c.BindJSON(&user)
// 		util.Check(err)

// 		if controller.SaveUser(user) {
// 			c.JSON(200, gin.H{})
// 		} else {
// 			c.Status(http.StatusBadRequest)
// 		}
// 		return
// 	}
// }

func login() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := model.UserInfo{}
		err := c.BindJSON(&user)
		util.Check(err)
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
	db := model.GetDBInstance()

	role := model.RoleInfo{}
	check(db.Where("id =?", roleID).Find(&role).Error)
	var j jsonPermission

	err := json.Unmarshal([]byte(role.Permission), &j)
	util.Check(err)

	v := reflect.ValueOf(&j).Elem().FieldByName(strings.Title(resource))
	arr := v.Interface().([]string)

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

func getRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		// AllowOrigins:    []strin g{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour}))

	r.Any("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	r.POST("/login", login())
	// r.POST("/register", register())

	r.GET("/users", controller.GetUsers())
	r.POST("/users", controller.AddUsers())               // add user
	r.PUT("/user/:user_name", controller.UpdateUser())    // update user
	r.DELETE("/user/:user_name", controller.DeleteUser()) // delete user

	r.GET("/projects", controller.GetProjects())
	r.POST("/projects", controller.AddProjects())                  // add project
	r.DELETE("/project/:project_name", controller.DeleteProject()) // delete project

	r.GET("/relations", controller.GetRelations())
	r.POST("/relations", controller.AddRelations())
	r.DELETE("/user/:user_name/project/:project_name/relation", controller.DeleteRelation())

	// r.POST("/user/:user_id/project", addProjectsForUser()) // assign project for user
	r.GET("/project/:project_name/envs", controller.GetProjectEnvs())     // add env for project
	r.POST("/project/:project_name/envs", controller.UpdateProjectEnvs()) // add env for project

	return r
}
