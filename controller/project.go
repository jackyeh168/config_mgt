package controller

import (
	"auth/model"
	"auth/util"

	"github.com/gin-gonic/gin"
)

func GetProjects() gin.HandlerFunc {
	return GetAll(&[]model.ProjectInfo{})
}

func SaveProject(project model.ProjectInfo) bool {
	return Save(&project)
}

func SaveProjects(projects []model.ProjectInfo) bool {
	res := true
	for _, project := range projects {
		res = SaveProject(project)
	}
	return res
}

func AddProjects() gin.HandlerFunc {
	return func(c *gin.Context) {
		var projects []model.ProjectInfo
		util.Check(c.BindJSON(&projects))

		SaveProjects(projects)
		c.JSON(200, gin.H{})
	}
}

func UpdateProject() gin.HandlerFunc {
	return func(c *gin.Context) {
		var project model.ProjectInfo
		util.Check(c.BindJSON(&project))
		// project.ID = util.StrToUint(c.Param("project_id"))

		db := model.GetDBInstance()
		var oldproject model.ProjectInfo
		oldproject.ID = util.StrToUint(c.Param("project_id"))
		db.First(&oldproject)

		oldproject.Name = project.Name

		db.Save(&oldproject)

		c.JSON(200, gin.H{})
		return
	}
}

func DeleteProject() gin.HandlerFunc {
	return func(c *gin.Context) {
		var project model.ProjectInfo
		project.ID = util.StrToUint(c.Param("project_id"))

		Delete(&project)

		c.JSON(200, gin.H{})
		return
	}
}
