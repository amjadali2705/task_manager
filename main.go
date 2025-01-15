package main

import (
	"task_manager/config"
	"task_manager/routes"
	"task_manager/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		utils.Logger.Fatal("Error loading .env file")
	}

	utils.InitLogger()
	defer utils.InitLogger()

	utils.Logger.Info("Starting the application...")

	config.InitDB()
	utils.Logger.Info("Database connection initialized")

	server := gin.Default()
	utils.Logger.Info("Server initialized")

	routes.RegisterRoutes(server)
	utils.Logger.Info("Routes registered")

	if err := server.Run(":8080"); err != nil {
		utils.Logger.Fatal("Failed to start the server", zap.Error(err))
	}
}
