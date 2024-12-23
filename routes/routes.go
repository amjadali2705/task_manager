package routes

import "github.com/gin-gonic/gin"

func RegisterRoutes(server *gin.Engine) {
	server.POST("/signup", signup)
	server.GET("/users", getUsers)
	server.GET("/logins", getLogins)
	server.GET("/tokens", getTokens)
}
