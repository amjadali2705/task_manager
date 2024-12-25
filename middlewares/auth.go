package middlewares

import (
	"net/http"
	"task_manager/utils"

	"github.com/gin-gonic/gin"
)

func Authenticate(context *gin.Context) {
	token := context.GetHeader("Authorization")
	if token == "" {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "token not found", "error": true})
		context.Abort()
		return
	}
	userId, err := utils.VerifyJwtToken(token)
	if err != nil {
		println(err)
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Not Authorized", "error": true})
		context.Abort()
		return
	}
	context.Set("userId", userId)
	context.Next()
}
