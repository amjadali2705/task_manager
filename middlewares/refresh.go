package middlewares

import (
	"net/http"
	"task_manager/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func decodeToken(token string) (int64, int64, error) {
	decoded, err := utils.DecodeJwtToken(token)
	if err != nil {
		return 0, 0, err
	}

	expiration := decoded["exp"].(float64)
	return int64(decoded["userId"].(float64)), int64(expiration), nil
}

func RefreshTokenMiddleware(context *gin.Context) {
	token := context.GetHeader("Authorization")
	if token == "" {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
		context.Abort()
		return
	}

	userId, expiration, err := decodeToken(token)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token", "error": true})
		context.Abort()
		return
	}

	now := time.Now().Unix()
	if expiration < now {
		refreshToken := context.GetHeader("Refresh-Token")
		if refreshToken == "" {
			context.JSON(401, gin.H{"message": "Refresh token is required", "error": true})
			context.Abort()
			return
		}

		newAccessToken, err := utils.GenerateJwtToken(userId)
		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"message": "Failed to generate new user token", "error": true})
			context.Abort()
			return
		}

		context.Header("New-Access-Token", newAccessToken)
	}

	context.Next()
}
