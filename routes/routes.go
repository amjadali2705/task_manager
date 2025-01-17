package routes

import (
	"task_manager/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.POST("/signUp", signUp)
	server.POST("/signIn", signIn)

	server.GET("avatar/:id", readAvatar)

	authenticated := server.Group("/")
	authenticated.Use(middlewares.Authenticate)
	authenticated.PUT("/users", updateUser)
	authenticated.GET("/tasks", getTasksByQuery)
	authenticated.GET("/tasks/:id", getTask)
	authenticated.GET("/users", getUser)
	authenticated.PUT("/updatePass", updatePassword)
	authenticated.POST("/tasks", createTask)
	authenticated.PUT("/tasks/:id", updateTask)
	authenticated.DELETE("/tasks/:id", deleteTask)
	authenticated.DELETE("signOut", signOut)
	authenticated.DELETE("signOut/all", signOutAll)

	authenticated.POST("avatar", uploadAvatar)
	authenticated.DELETE("avatar", deleteAvatar)

	server.POST("/refresh-token", refreshTokenHandler)
}
