package middlewares

import (
	"fmt"
	"net/http"
	"strings"
	"task_manager/config"
	"task_manager/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Authenticate(context *gin.Context) {
	token := context.Request.Header.Get("Authorization")

	if token == "" {
		utils.Logger.Warn("Authorization token is missing", zap.String("method", context.Request.Method), zap.String("url", context.Request.URL.String()))
		context.JSON(http.StatusUnauthorized, gin.H{"message": "token not found", "error": true})
		context.Abort()
		return
	}

	token = strings.TrimPrefix(token, "Bearer ")

	userId, err := utils.VerifyJwtToken(token)
	if err != nil {
		utils.Logger.Error("Failed to verify user token", zap.Error(err))
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
		utils.Logger.Error("Session expired or token not found", zap.Error(err))
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Session Expired: User has to log in", "error": true})
	}

	utils.Logger.Info("Token found in the database", zap.String("tokenId", fmt.Sprintf("%d", dbToken.ID)))
	return err
}
