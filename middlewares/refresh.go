package middlewares

import (
	"fmt"
	"net/http"
	"task_manager/utils"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func decodeToken(token string) (int64, int64, error) {
	decoded, err := utils.DecodeJwtToken(token)
	if err != nil {
		utils.Logger.Error("Failed to decode user token", zap.Error(err))
		return 0, 0, err
	}

	expiration := decoded["exp"].(float64)
	userId := int64(decoded["userId"].(float64))

	utils.Logger.Info("User Token decoded successfully", zap.String("userId", fmt.Sprintf("%d", userId)), zap.Int64("expiration", int64(expiration)))

	return userId, int64(expiration), nil
}

func RefreshTokenMiddleware(context *gin.Context) {
	token := context.GetHeader("Authorization")

	if token == "" {
		utils.Logger.Warn("Authorization token is missing", zap.String("method", context.Request.Method), zap.String("url", context.Request.URL.String()))
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
		context.Abort()
		return
	}

	userId, expiration, err := decodeToken(token)
	if err != nil {
		utils.Logger.Error("Invalid token during refresh process", zap.Error(err))
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token", "error": true})
		context.Abort()
		return
	}

	now := time.Now().Unix()
	if expiration < now {
		utils.Logger.Warn("Token has expired", zap.Int64("userId", userId))

		refreshToken := context.GetHeader("Refresh-Token")
		if refreshToken == "" {
			utils.Logger.Warn("Refresh token is required", zap.String("method", context.Request.Method), zap.String("url", context.Request.URL.String()))
			context.JSON(401, gin.H{"message": "Refresh token is required", "error": true})
			context.Abort()
			return
		}

		newUserToken, err := utils.GenerateJwtToken(userId)
		if err != nil {
			utils.Logger.Error("Failed to generate new user token", zap.Error(err), zap.String("userId", fmt.Sprintf("%d", userId)))
			context.JSON(http.StatusUnauthorized, gin.H{"message": "Failed to generate new user token", "error": true})
			context.Abort()
			return
		}

		utils.Logger.Info("New user token generated successfully")

		context.Header("New-User-Token", newUserToken)
	}

	context.Next()
}
