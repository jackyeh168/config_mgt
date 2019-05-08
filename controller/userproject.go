package controller

import (
	"auth/model"
	"auth/util"

	"github.com/gin-gonic/gin"
)

func GetUserProjects() gin.HandlerFunc {
	return GetAll(&[]model.UserProject{})
}

func SaveUserProject(userProject model.UserProject) bool {
	return Save(&userProject)
}

func SaveUserProjects(userProjects []model.UserProject) bool {
	res := true
	for _, userProject := range userProjects {
		res = SaveUserProject(userProject)
	}
	return res
}

func AddUserProjects() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userProjects []model.UserProject
		util.Check(c.BindJSON(&userProjects))

		SaveUserProjects(userProjects)
		c.JSON(200, gin.H{})
	}
}

func DeleteUserProject() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userProject model.UserProject
		util.Check(c.BindJSON(&userProject))

		Delete(userProject)
		c.JSON(200, gin.H{})
	}
}
