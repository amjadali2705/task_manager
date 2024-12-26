package middlewares

import (
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
		context.JSON(401, gin.H{"error": "Authorization token is required"})
		context.Abort()
		return
	}

	userId, expiration, err := decodeToken(token)
	if err != nil {
		context.JSON(401, gin.H{"error": "Invalid token"})
		context.Abort()
		return
	}

	now := time.Now().Unix()
	if expiration < now {
		refreshToken := context.GetHeader("Refresh-Token")
		if refreshToken == "" {
			context.JSON(401, gin.H{"error": "Refresh token is required"})
			context.Abort()
			return
		}

		newAccessToken, err := utils.GenerateJwtToken(userId)
		if err != nil {
			context.JSON(500, gin.H{"error": "Internal server error"})
			context.Abort()
			return
		}

		context.Header("New-Access-Token", newAccessToken)
	}

	context.Next()
}
