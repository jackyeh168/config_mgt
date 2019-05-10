package controller

import (
	"auth/dynamo"
	"auth/util"

	"github.com/gin-gonic/gin"
)

func GetProjects() gin.HandlerFunc {
	return func(c *gin.Context) {
		err, projects := dynamo.GetProjects()
		util.Check(err)

		c.JSON(200, gin.H{
			"data": projects,
		})
	}
}

func AddProjects() gin.HandlerFunc {
	return func(c *gin.Context) {
		var projects []dynamo.Project
		util.Check(c.BindJSON(&projects))

		dynamo.SaveProjects(projects)
		c.JSON(200, gin.H{})
	}
}

func DeleteProject() gin.HandlerFunc {
	return func(c *gin.Context) {
		project := dynamo.Project{ProjectName: c.Param("project_name")}
		dynamo.DeleteProject(project)

		c.JSON(200, gin.H{})
	}
}

func GetProjectEnvs() gin.HandlerFunc {
	return func(c *gin.Context) {
		err, projectenvs := dynamo.GetProjectEnvs(c.Param("project_name"))
		util.Check(err)

		c.JSON(200, gin.H{
			"data": projectenvs,
		})
	}
}

func UpdateProjectEnvs() gin.HandlerFunc {
	return func(c *gin.Context) {
		var project dynamo.Project
		util.Check(c.BindJSON(&project))

		dynamo.UpdateProjectEnvs(project)
		c.JSON(200, gin.H{})
	}
}
