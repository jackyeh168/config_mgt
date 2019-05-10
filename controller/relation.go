package controller

import (
	"auth/dynamo"
	"auth/util"

	"github.com/gin-gonic/gin"
)

func GetRelations() gin.HandlerFunc {
	return func(c *gin.Context) {
		err, users := dynamo.GetRelations()
		util.Check(err)

		c.JSON(200, gin.H{
			"data": users,
		})
	}
}

func AddRelations() gin.HandlerFunc {
	return func(c *gin.Context) {
		var relations []dynamo.Relation
		util.Check(c.BindJSON(&relations))

		dynamo.SaveRelations(relations)
		c.JSON(200, gin.H{})
	}
}

func DeleteRelation() gin.HandlerFunc {
	return func(c *gin.Context) {
		relation := dynamo.Relation{
			UserName:    c.Param("user_name"),
			ProjectName: c.Param("project_name"),
		}

		dynamo.DeleteRelation(relation)
		c.JSON(200, gin.H{})
	}
}
