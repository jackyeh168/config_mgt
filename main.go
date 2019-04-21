package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func getHTTPPort() (port string) {
	port = os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	return
}

func main() {
	r := gin.Default()
	r.Use(gin.Logger())
	// r.Use(gin.Recovery())

	r.POST("/login", login())

	authorized := r.Group("/")
	// per group middleware! in this case we use the custom created
	// AuthRequired() middleware just in the "authorized" group.
	authorized.Use(authRequired())
	{
		authorized.GET("/getData", func(c *gin.Context) {

			c.JSON(200, gin.H{
				"data": "justfortest",
			})
		})
	}

	r.Run(":" + getHTTPPort()) // listen and serve on 0.0.0.0:8080
}
