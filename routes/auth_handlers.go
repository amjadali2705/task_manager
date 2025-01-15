// package routes

// import (
// 	"net/http"
// 	"strings"
// 	"task_manager/models"
// 	"task_manager/utils"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"go.uber.org/zap"
// )

// func refreshTokenhandler(c *gin.Context) {
// 	var user models.User

// 	refreshToken := c.GetHeader("Refresh-Token")
// 	if refreshToken == "" {
// 		utils.Logger.Warn("Missing Refresh Token in request header", zap.String("url", c.Request.URL.String()))
// 		c.JSON(http.StatusUnauthorized, gin.H{"message": "Refresh token is required", "error": true})
// 		c.Abort()
// 		return
// 	}

// 	userId, err := utils.VerifyRefreshToken(refreshToken)
// 	if err != nil {
// 		utils.Logger.Error("Invalid refresh token", zap.Error(err))
// 		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid refresh token", "error": true})
// 		c.Abort()
// 		return
// 	}

// 	userToken := c.GetHeader("Authorization")
// 	if userToken == "" {
// 		utils.Logger.Warn("Missing User Token in request header", zap.String("url", c.Request.URL.String()))
// 		c.JSON(http.StatusUnauthorized, gin.H{"message": "User token is required", "error": true})
// 		c.Abort()
// 		return
// 	}

// 	userToken = strings.TrimPrefix(userToken, "Bearer ")

// 	// Decode the user token to check its expiration
// 	claims, err := utils.DecodeJwtToken(userToken)
// 	if err != nil {
// 		utils.Logger.Warn("Invalid User Token", zap.Error(err))
// 		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid user token", "error": true})
// 		c.Abort()
// 		return
// 	}

// 	// Check if the token has expired
// 	exp, ok := claims["exp"].(float64) // `exp` is usually a Unix timestamp
// 	if !ok {
// 		utils.Logger.Error("User token does not contain a valid expiration claim")
// 		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token claims", "error": true})
// 		c.Abort()
// 		return
// 	}

// 	if int64(exp) > time.Now().Unix() {
// 		// Token is still valid, no need to generate new tokens
// 		utils.Logger.Info("User token is still valid", zap.Int64("userId", userId))
// 		c.JSON(http.StatusOK, gin.H{"message": "User token is still valid", "error": false})
// 		return
// 	}

// 	// Generate new tokens only if the user token has expired
// 	newUserToken, err := utils.GenerateJwtToken(userId)
// 	if err != nil {
// 		utils.Logger.Error("Error generating new user token", zap.Error(err), zap.Int64("userId", userId))
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error generating new user token", "error": true})
// 		c.Abort()
// 		return
// 	}

// 	newRefreshToken, err := utils.GenerateRefreshToken(userId)
// 	if err != nil {
// 		utils.Logger.Error("Error generating new refresh token", zap.Error(err), zap.Int64("userId", userId))
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error generating new refresh token", "error": true})
// 		c.Abort()
// 		return
// 	}

// 	err = user.SaveToken(newUserToken, newRefreshToken)
// 	if err != nil {
// 		utils.Logger.Error("Error saving new tokens to database", zap.Error(err), zap.Int64("userId", userId))
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving new tokens", "error": true})
// 		c.Abort()
// 		return
// 	}

// 	err = models.DeleteRefreshToken(refreshToken)
// 	if err != nil {
// 		utils.Logger.Warn("Error deleting old refresh token", zap.Error(err))
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error deleting old refresh token", "error": true})
// 		c.Abort()
// 		return
// 	}

// 	utils.Logger.Info("Tokens refreshed successfully", zap.Int64("userId", userId))
// 	c.JSON(http.StatusOK, gin.H{
// 		"message":       "Tokens refreshed successfully",
// 		"error":         false,
// 		"user_token":    newUserToken,
// 		"refresh_token": newRefreshToken,
// 	})
// }

package routes

import (
	"net/http"
	"task_manager/middlewares"
	"task_manager/models"
	"task_manager/utils"

	"github.com/gin-gonic/gin"
)

// RefreshTokenHandler handles the refresh token logic
func refreshTokenHandler(c *gin.Context) {
	// Get the refresh token from the request header
	var user models.User

	err := middlewares.CheckRefreshToken(c)
	if err != nil {
		return
	}

	refreshToken := c.GetHeader("Refresh-Token")
	if refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Refresh token required", "error": true, "data": nil})
		return
	}

	// Verify the refresh token
	userId, err := utils.VerifyRefreshToken(refreshToken)
	if err != nil {
		err = models.DeleteRefreshToken(refreshToken)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Failed to delete refresh token", "error": true, "data": nil})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid refresh token", "error": true, "data": nil})
		return
	}

	// Generate a new access token
	newUserToken, err := utils.GenerateJwtToken(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error generating new access token", "error": true, "data": nil})
		return
	}

	newRefreshToken, err := utils.GenerateRefreshToken(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error generating new refresh token", "error": true, "data": nil})
		return
	}

	err = user.SaveToken(newUserToken, newRefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not save the token", "error": true, "data": nil})
		return
	}

	err = models.DeleteRefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Failed to delete refresh token", "error": true, "data": nil})
		return
	}

	// Return the new access token to the client
	c.JSON(http.StatusOK, gin.H{"message": "Token refreshed successfully", "data": gin.H{"user_token": newUserToken, "refresh_token": newRefreshToken}, "error": false})
}
