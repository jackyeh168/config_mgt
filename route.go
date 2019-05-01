package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const hmacSampleSecret = "GYthbtJJ6tp3852JMEVmVHhDckdHHDsJ"

type jsonPermission struct {
	User    []string `json: "user"`
	Project []string `json: "project"`
	Env     []string `json: "env"`
}

func strToUint(s string) uint {
	i, err := strconv.Atoi(s)
	check(err)
	return (uint)(i)
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

func saveUser(user UserInfo) bool {
	db := getDBInstance()

	if db.Where("username = ?", user.Username).Find(&user).RecordNotFound() {
		user.Password = encrypt(user.Password)
		check(db.Create(&user).Error)
		return true
	}
	return false
}

func saveUsers(users []UserInfo) bool {
	res := true
	for _, user := range users {
		res = saveUser(user)
	}
	return res
}

func addUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []UserInfo
		c.BindJSON(&users)

		saveUsers(users)
		c.JSON(200, gin.H{})
	}
}

func getProjects() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := getDBInstance()
		var projects []ProjectInfo
		db.Find(&projects)

		c.JSON(200, gin.H{
			"data": projects,
		})
	}
}

func saveProject(project ProjectInfo) bool {
	db := getDBInstance()

	if db.Where("name = ?", project.Name).Find(&project).RecordNotFound() {
		check(db.Create(&project).Error)
		return true
	}
	return false
}

func saveProjects(projects []ProjectInfo) bool {
	res := true
	for _, project := range projects {
		res = saveProject(project)
	}
	return res
}

func addProjects() gin.HandlerFunc {
	return func(c *gin.Context) {
		var projects []ProjectInfo
		c.BindJSON(&projects)

		saveProjects(projects)
		c.JSON(200, gin.H{})
	}
}

func getProjectEnvs() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := getDBInstance()
		projectID := c.Param("project_id")
		var envs []ProjectEnv
		db.Where("project_id = ?", projectID).Find(&envs)

		c.JSON(200, gin.H{
			"data": envs,
		})
	}
}

func saveEnv(env ProjectEnv) bool {
	if env.EnvKey == "" {
		return false
	}

	db := getDBInstance()
	var oldEnv ProjectEnv
	newdb := db.Where("project_id = ? and env_key = ?", env.ProjectID, env.EnvKey).Find(&oldEnv)
	if newdb.RecordNotFound() {
		check(db.Create(&env).Error)
	} else {
		check(newdb.Save(&env).Error)
	}
	return true
}

func saveEnvs(project_id uint, envs []ProjectEnv) bool {
	res := true
	for _, env := range envs {
		if env.IsSecret {
			env.EnvValue = encrypt(env.EnvValue)
		}
		env.ProjectID = project_id
		res = saveEnv(env)
	}
	return res
}

func deleteAllEnvs(project_id uint) {
	db := getDBInstance()
	check(db.Where("project_id = ?", project_id).Unscoped().Delete(ProjectEnv{}).Error)
}

func addProjectEnvs() gin.HandlerFunc {
	return func(c *gin.Context) {
		var projectEnvs []ProjectEnv
		c.BindJSON(&projectEnvs)

		project_id := strToUint(c.Param("project_id"))
		deleteAllEnvs(project_id)
		saveEnvs(project_id, projectEnvs)
		c.JSON(200, gin.H{})
	}
}

func getUserProjects() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := getDBInstance()
		userID := c.Param("user_id")
		userProjects := []UserProject{}
		db.Where("user_id = ?", userID).Find(&userProjects)

		c.JSON(200, gin.H{
			"data": userProjects,
		})
	}
}

func saveUserProject(userProject UserProject) bool {
	db := getDBInstance()
	err := db.Where(&userProject).Find(&userProject)
	if err.RecordNotFound() {
		check(db.Create(&userProject).Error)
	} else {
		check(db.Save(&userProject).Error)
	}
	return true
}

func saveUserProjects(userProjects []UserProject) bool {
	res := true
	for _, userProject := range userProjects {
		res = saveUserProject(userProject)
	}
	return res
}

func addUserProject() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userProjects []UserProject
		c.BindJSON(&userProjects)

		saveUserProjects(userProjects)
		c.JSON(200, gin.H{})
	}
}

func getRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())

	r.POST("/login", login())
	r.POST("/register", register())

	r.GET("/users", getUsers())
	r.POST("/users", addUsers()) // add user
	// r.POST("/user/:user_id/project", addProjectsForUser()) // assign project for user
	r.GET("/project/:project_id/envs", getProjectEnvs())  // add env for project
	r.POST("/project/:project_id/envs", addProjectEnvs()) // add env for project
	r.GET("/user/:user_id/projectowner", getUserProjects())
	r.POST("/projectowner", addUserProject())
	r.GET("/projects", getProjects())
	r.POST("/projects", addProjects()) // add project

	// authorized := r.Group("/")
	// // per group middleware! in this case we use the custom created
	// // AuthRequired() middleware just in the "authorized" group.
	// authorized.Use(authRequired())
	// {
	// 	authorized.GET("/users", isValid("user", "read"), getUsers())
	// 	authorized.POST("/users", isValid("user", "create"), addUsers()) // add user
	// 	// authorized.POST("/user/:user_id/project", isValid("project", "create"), addProjectsForUser()) // assign project for user

	// 	authorized.GET("/project/:project_id/envs", getProjectEnvs()) // add env for project
	// 	authorized.POST("/project/:project_id/env", addProjectEnvs()) // add env for project

	// 	// authorized.GET("/projectowners", getAllUserProjects())
	// 	authorized.GET("/user/:user_id/projectowner", getUserProjects())
	// 	authorized.POST("/projectowner", addUserProject())

	// 	authorized.GET("/projects", getProjects())
	// 	authorized.POST("/projects", addProjects()) // add project
	// }
	return r
}
