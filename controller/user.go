package controller

import (
	"auth/dynamo"
	"auth/util"

	"github.com/gin-gonic/gin"
)

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		err, users := dynamo.GetUsers()
		util.Check(err)

		c.JSON(200, gin.H{
			"data": users,
		})
	}
}

func AddUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []dynamo.User
		util.Check(c.BindJSON(&users))

		dynamo.SaveUsers(users)
		c.JSON(200, gin.H{})
	}
}

func UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user dynamo.User
		util.Check(c.BindJSON(&user))

		dynamo.UpdateUser(user)

		c.JSON(200, gin.H{})
		return
	}
}

func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := dynamo.User{UserName: c.Param("user_name")}
		util.Check(dynamo.DeleteUser(user))

		c.JSON(200, gin.H{})
		return
	}
}
