package routes

import (
	"net/http"
	"task_manager/models"
	"task_manager/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RefreshTokenhandler(c *gin.Context) {
	var user models.User

	refreshToken := c.GetHeader("Refresh-Token")
	if refreshToken == "" {
		utils.Logger.Warn("Missing Refresh Token in request header", zap.String("url", c.Request.URL.String()))
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Refresh token is required", "error": true})
		c.Abort()
		return
	}

	userId, err := utils.VerifyRefreshToken(refreshToken)
	if err != nil {
		utils.Logger.Error("Invalid refresh token", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid refresh token", "error": true})
		c.Abort()
		return
	}

	newUserToken, err := utils.GenerateJwtToken(userId)
	if err != nil {
		utils.Logger.Error("Error generating new user token", zap.Error(err), zap.Int64("userId", userId))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error generating new user token", "error": true})
		c.Abort()
		return
	}

	newRefreshToken, err := utils.GenerateRefreshToken(userId)
	if err != nil {
		utils.Logger.Error("Error generating new refresh token", zap.Error(err), zap.Int64("userId", userId))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error generating new refresh token", "error": true})
		c.Abort()
		return
	}

	err = user.SaveToken(newUserToken, newRefreshToken)
	if err != nil {
		utils.Logger.Error("Error saving new tokens to database", zap.Error(err), zap.Int64("userId", userId))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving new tokens", "error": true})
		c.Abort()
		return
	}

	err = models.DeleteRefreshToken(refreshToken)
	if err != nil {
		utils.Logger.Warn("Error deleting old refresh token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error deleting old refresh token", "error": true})
		c.Abort()
		return
	}

	utils.Logger.Info("Tokens refreshed successfully", zap.Int64("userId", userId))
	c.JSON(http.StatusOK, gin.H{"message": "Tokens refreshed successfully", "error": false, "user_token": newUserToken, "refresh_token": newRefreshToken})
}
