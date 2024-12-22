package main

import (
	"task_manager/db"
	"task_manager/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	server := gin.Default()

	server.POST("/signup", routes.Signup)

	server.Run(":8080")
}
