package controller

import (
	"auth/model"
	"auth/util"

	"github.com/gin-gonic/gin"
)

func GetProjectEnvs() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := model.GetDBInstance()
		projectID := c.Param("project_id")
		var envs []model.ProjectEnv
		db.Where("project_id = ?", projectID).Find(&envs)

		c.JSON(200, gin.H{
			"data": envs,
		})
	}
}

func SaveEnv(env model.ProjectEnv) bool {
	db := model.GetDBInstance()
	var oldEnv model.ProjectEnv
	newdb := db.Where("project_id = ? and env_key = ?", env.ProjectID, env.EnvKey).Find(&oldEnv)
	if newdb.RecordNotFound() {
		util.Check(db.Create(&env).Error)
	} else {
		util.Check(newdb.Save(&env).Error)
	}
	return true
}

func SaveEnvs(project_id uint, envs []model.ProjectEnv) bool {
	res := true
	for _, env := range envs {
		if len(env.EnvValue) < 2 || env.EnvValue[:2] != "$2" {
			env.EnvValue = util.Encrypt(env.EnvValue)
		}
		env.ProjectID = project_id
		res = SaveEnv(env)
	}
	return res
}

func DeleteAllEnvs(project_id uint) {
	db := model.GetDBInstance()
	util.Check(db.Where("project_id = ?", project_id).Unscoped().Delete(model.ProjectEnv{}).Error)
}

func AddProjectEnvs() gin.HandlerFunc {
	return func(c *gin.Context) {
		var projectEnvs []model.ProjectEnv
		c.BindJSON(&projectEnvs)

		project_id := util.StrToUint(c.Param("project_id"))
		DeleteAllEnvs(project_id)
		SaveEnvs(project_id, projectEnvs)
		c.JSON(200, gin.H{})
	}
}
