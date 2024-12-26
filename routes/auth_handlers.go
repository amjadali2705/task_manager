package routes

import (
	"net/http"
	"task_manager/utils"

	"github.com/gin-gonic/gin"
)

func RefreshTokenhandler(c *gin.Context) {
	refreshToken := c.GetHeader("Refresh-Token")
	if refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Refresh token is required", "error": true})
		c.Abort()
		return
	}

	userId, err := utils.VerifyRefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid refresh token", "error": true})
		c.Abort()
		return
	}

	newAccessToken, err := utils.GenerateJwtToken(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error generating new access token", "error": true})
		c.Abort()
		return
	}

	newRefreshToken, err := utils.GenerateRefreshToken(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error generating new refresh token", "error": true})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tokens refreshed successfully", "error": false, "user_token": newAccessToken, "refresh_token": newRefreshToken})
}
