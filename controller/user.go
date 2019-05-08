package controller

import (
	"auth/model"
	"auth/util"

	"github.com/gin-gonic/gin"
)

func GetUsers() gin.HandlerFunc {
	return GetAll(&[]model.UserInfo{})
}

func SaveUser(user model.UserInfo) bool {
	user.Password = util.Encrypt(user.Password)
	return Save(&user)
}

func SaveUsers(users []model.UserInfo) bool {
	res := true
	for _, user := range users {
		res = SaveUser(user)
	}
	return res
}

func AddUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []model.UserInfo
		util.Check(c.BindJSON(&users))
		SaveUsers(users)
		c.JSON(200, gin.H{})
	}
}

func UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user model.UserInfo
		util.Check(c.BindJSON(&user))

		db := model.GetDBInstance()
		var olduser model.UserInfo
		olduser.ID = util.StrToUint(c.Param("user_id"))
		db.First(&olduser)

		if user.Password != "" {
			olduser.Password = util.Encrypt(user.Password)
		}
		olduser.Username = user.Username
		olduser.RoleID = user.RoleID
		db.Save(&olduser)

		c.JSON(200, gin.H{})

		return
	}
}

func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user model.UserInfo
		user.ID = util.StrToUint(c.Param("user_id"))

		db := model.GetDBInstance()
		db.Unscoped().Delete(&user)

		c.JSON(200, gin.H{})
		return
	}
}
