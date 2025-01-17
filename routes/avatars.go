package routes

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"task_manager/middlewares"
	"task_manager/models"
	"task_manager/utils"

	"github.com/gin-gonic/gin"
)

func uploadAvatar(c *gin.Context) {
	err := middlewares.CheckTokenPresent(c)
	if err != nil {
		return
	}

	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not authenticated", "data": nil, "error": true})
		// utils.StandardResponse(c, http.StatusUnauthorized, "User not authenticated", true, nil)
		return
	}

	_, err = models.ReadAvatar(userId.(int64))
	if err == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Avatar already created, please update avatar", "data": nil, "error": true})
		// utils.StandardResponse(c, http.StatusBadRequest, "Avatar already created, please update avatar", true, nil)
		return
	}

	fileHeader, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid file", "data": nil, "error": true})
		// utils.StandardResponse(c, http.StatusBadRequest, "Invalid file", true, nil)
		return
	}

	fileExtension, err := utils.ValidateAvatar(fileHeader)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "error": true, "data": nil})
		// utils.StandardResponse(c, http.StatusBadRequest, err.Error(), true, nil)
		return
	}

	fileName := fmt.Sprintf("avatar_%s%s", strconv.FormatInt(userId.(int64), 10), fileExtension)

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to open uploaded file", "data": nil, "error": true})
		// utils.StandardResponse(c, http.StatusInternalServerError, "Failed to open uploaded file", true, nil)
		return
	}

	content, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to read uploded file", "error": true, "data": nil})
		// utils.StandardResponse(c, http.StatusInternalServerError, "Failed to read uploaded file", true, nil)
		return
	}

	err = models.SaveAvatar(userId.(int64), content, fileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save avatar", "error": true, "data": nil})
		// utils.StandardResponse(c, http.StatusInternalServerError, "Failed to save avatar", true, nil)
		return
	}

	defer file.Close()

	c.JSON(http.StatusOK, gin.H{"message": "Avatar uploded", "error": false, "data": nil})
	// utils.StandardResponse(c, http.StatusOK, "Avatar uploded successfully", false, nil)

}
