package middlewares

import (
	"net/http"
	"strings"
	"task_manager/config"
	"task_manager/utils"

	"github.com/gin-gonic/gin"
)

func Authenticate(context *gin.Context) {
	token := context.Request.Header.Get("Authorization")
	if token == "" {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "token not found", "error": true})
		context.Abort()
		return
	}

	token = strings.TrimPrefix(token, "Bearer ")

	userId, err := utils.VerifyJwtToken(token)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Not Authorized", "error": true})
		context.Abort()
		return
	}

	context.Set("token", token)
	context.Set("userId", userId)

	context.Next()
}

func CheckTokenPresent(context *gin.Context) error {
	token := context.Request.Header.Get("Authorization")

	token = strings.TrimPrefix(token, "Bearer ")

	var dbToken config.Token

	err := config.DB.Where("user_token = ?", token).First(&dbToken).Error
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Session Expired: User has to log in", "error": true})
	}
	return err
}
