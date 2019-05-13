package main

import (
	"auth/controller"
	"auth/dynamo"
	"auth/util"
	"net/http"
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

func verifyUser(user dynamo.User) (bool, dynamo.User) {
	err, savedUser := dynamo.GetUser(user.UserName)
	util.Check(err)

	if savedUser.UserName == "" {
		// the user doesn't exist
		return false, dynamo.User{}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(savedUser.Password), []byte(user.Password)); err != nil {
		return false, dynamo.User{}
	}

	return true, savedUser
}

func getUserRelationList(userName string) []string {

	err, relations := dynamo.GetUserRelations(userName)
	util.Check(err)

	listLen := len(relations)
	relationList := make([]string, listLen)

	for i := 0; i < listLen; i++ {
		relationList[i] = relations[i].ProjectName
	}

	return relationList
}

func stringifyStrArr(s []string) string {
	return strings.Join(s, ",")
}

func login() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := dynamo.User{}
		err := c.BindJSON(&user)
		util.Check(err)

		if isValid, savedUser := verifyUser(user); isValid {
			relation := stringifyStrArr(getUserRelationList(user.UserName))
			c.JSON(200, gin.H{
				"token": getJWT(savedUser, relation),
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
		if err, roleName, relation := verifyJWT(token); err == nil {
			// if success, then next
			c.Set("roleName", roleName)
			c.Set("relation", relation)
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

func verifyProjectPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verify role
		roleName, roleNameOK := c.Get("roleName")
		allowedRelation, relationOK := c.Get("relation")
		expectRelation := c.Param("project_name")

		if !roleNameOK || !relationOK {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if roleName == "admin" {
			// if permission is valid, then next
			c.Next()
			return
		} else if roleName == "user" && strings.Contains(allowedRelation.(string), expectRelation) {
			c.Next()
			return
		}
		// else return http 401
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}

func isAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verify role
		roleName, ok := c.Get("roleName")
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if roleName == "admin" {
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

	authorized := r.Group("/")
	authorized.Use(authRequired())
	{
		authorized.GET("/projects", controller.GetProjects())

		adminRoute := authorized.Group("/")
		adminRoute.Use(isAdmin())
		{
			// only admin can access user info
			adminRoute.GET("/users", controller.GetUsers())
			adminRoute.POST("/users", controller.AddUsers())               // add user
			adminRoute.PUT("/user/:user_name", controller.UpdateUser())    // update user
			adminRoute.DELETE("/user/:user_name", controller.DeleteUser()) // delete user

			// only admin can access projects
			adminRoute.POST("/projects", controller.AddProjects())                  // add project
			adminRoute.DELETE("/project/:project_name", controller.DeleteProject()) // delete project

			// only admin can access relations
			adminRoute.GET("/relations", controller.GetRelations())
			adminRoute.POST("/relations", controller.AddRelations())
			adminRoute.DELETE("/user/:user_name/project/:project_name/relation", controller.DeleteRelation())
		}

		// users can access all projects' envs if the role is admin,
		// otherwise, users can only access projects they rely on
		authorized.GET("/project/:project_name/envs", verifyProjectPermission(), controller.GetProjectEnvs())     // get env for project
		authorized.POST("/project/:project_name/envs", verifyProjectPermission(), controller.UpdateProjectEnvs()) // add env for project
	}
	return r
}
