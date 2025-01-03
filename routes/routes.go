package routes

import (
	"task_manager/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.POST("/signUp", signUp)
	server.POST("/signIn", signIn)

	// server.GET("/tasks", getTasks)
	// server.GET("/tasks", getTasksByQuery)

	// server.GET("/tasks/:id", getTask)
	// server.GET("/users", getUsers)
	// server.GET("/logins", getLogins)
	// server.GET("/tokens", getTokens)

	authenticated := server.Group("/")
	authenticated.Use(middlewares.Authenticate)
	authenticated.PUT("/users", updateUser)
	authenticated.GET("/tasks", getTasksByQuery)
	authenticated.GET("/tasks/:id", getTask)
	authenticated.POST("/tasks", createTask)
	authenticated.PUT("/tasks/:id", updateTask)
	authenticated.DELETE("/tasks/:id", deleteTask)
	authenticated.DELETE("signOut", signOut)

	server.POST("/refresh-token", RefreshTokenhandler)
}
