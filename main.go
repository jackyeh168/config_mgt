package main

import (
	"os"
)

func getHTTPPort() (port string) {
	port = os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	return
}

func launchServer() {
	r := getRouter()
	r.Run(":" + getHTTPPort()) // listen and serve on 0.0.0.0:8080
}

func main() {
	initDB()
	// migration()
	// seed()

	launchServer()
}
