package routes

import (
	"task_manager/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.POST("/signup", signup)
	server.POST("/login", login)

	
	server.GET("/tasks", getTasks)
	server.GET("/tasks/:id", getTask)
	server.GET("/users", getUsers)
	// server.GET("/logins", getLogins)
	// server.GET("/tokens", getTokens)

	authenticated := server.Group("/")
	authenticated.Use(middlewares.Authenticate)
	authenticated.POST("/tasks", createTask)
	authenticated.PUT("/tasks/:id", updateTask)
	authenticated.DELETE("/tasks/:id", deleteTask)
}

