package main

import (
	"task_manager/config"
	"task_manager/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitDB()
	server := gin.Default()

	routes.RegisterRoutes(server)

	server.Run(":8080")
}
