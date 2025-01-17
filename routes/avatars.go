package routes

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"task_manager/middlewares"
	"task_manager/models"
	"task_manager/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func uploadAvatar(c *gin.Context) {
	// Check if the token is present and valid
	err := middlewares.CheckTokenPresent(c)
	if err != nil {
		return
	}

	// Retrieve the userId from the context
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not authenticated", "data": nil, "error": true})
		return
	}

	// Check if an avatar already exists for the user
	existingAvatar, err := models.ReadAvatar(userId.(int64))
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to check avatar existence", "data": nil, "error": true})
		return
	}

	// Handle the uploaded file
	fileHeader, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid file", "data": nil, "error": true})
		return
	}

	// Validate the file (e.g., check size, type)
	fileExtension, err := utils.ValidateAvatar(fileHeader)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "data": nil, "error": true})
		return
	}

	fileName := fmt.Sprintf("avatar_%s%s", strconv.FormatInt(userId.(int64), 10), fileExtension)

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to open uploaded file", "data": nil, "error": true})
		return
	}
	defer file.Close()

	// Read the file content
	content, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to read uploaded file", "data": nil, "error": true})
		return
	}

	// Save or update the avatar in the database
	if existingAvatar != nil {
		// Update the existing avatar
		err = models.UpdateAvatar(userId.(int64), content, fileName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update avatar", "data": nil, "error": true})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Avatar updated successfully", "data": nil, "error": false})
	} else {
		// Create a new avatar
		err = models.SaveAvatar(userId.(int64), content, fileName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save avatar", "data": nil, "error": true})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Avatar uploaded successfully", "data": nil, "error": false})
	}
}

func readAvatar(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Logger.Error("Failed to parse avatar Id", zap.String("param", c.Param("id")), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse avatar id", "error": true, "data": nil})
		// utils.StandardResponse(c, http.StatusBadRequest, "Could not parse avatar id", true, nil)
		return
	}

	avatar, err := models.ReadAvatar(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch avatar", "error": true})
		// utils.StandardResponse(c, http.StatusInternalServerError, "Could not fetch avatar", true, nil)
		return
	}

	c.Data(http.StatusOK, "image/jpg", avatar.Data)
	// utils.StandardResponse(c, http.StatusOK, "Avatar fetched successfully", true, avatar.Data)
}

func deleteAvatar(c *gin.Context) {
	err := middlewares.CheckTokenPresent(c)
	if err != nil {
		return
	}

	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized: User not authenticated", "error": true, "data": nil})
		// utils.StandardResponse(c, http.StatusUnauthorized, "Unauthorized, User not authenticated", true, nil)
		return
	}

	_, err = models.ReadAvatar(userId.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "No avatar present to delete", "data": nil, "error": true})
		// utils.StandardResponse(c, http.StatusInternalServerError, "No avatar present to delete", true, nil)
		return
	}

	err = models.DeleteAvatar(userId.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete avatar", "error": true, "data": nil})
		// utils.StandardResponse(c, http.StatusInternalServerError, "Failed to delete avatar", true, nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Avatar deleted", "error": false, "data": nil})
	// utils.StandardResponse(c, http.StatusOK, "Avatar deleted successfully", false, nil)
}
