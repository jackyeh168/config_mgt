package controller

import (
	"auth/model"
	"auth/util"

	"github.com/gin-gonic/gin"
)

func GetAll(i interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := model.GetDBInstance()
		db.Find(i)

		c.JSON(200, gin.H{
			"data": i,
		})
	}
}

func Get(i []interface{}, id uint) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := model.GetDBInstance()
		db.Where("ID = ?", id).Find(i)

		c.JSON(200, gin.H{
			"data": i,
		})
	}
}

func Save(i interface{}) bool {
	db := model.GetDBInstance()

	if db.Where(i).Find(i).RecordNotFound() {
		util.Check(db.Create(i).Error)
		return true
	}
	return false
}

func Delete(i interface{}) bool {
	db := model.GetDBInstance()

	if err := db.Where(i).Unscoped().Delete(i).Error; err != nil {
		return false
	}
	return true
}
