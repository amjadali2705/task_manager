package routes

import (
	"net/http"
	"strings"
	"task_manager/middlewares"
	"task_manager/models"
	"task_manager/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func signUp(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)
	if err != nil {
		utils.Logger.Error("Failed to parse request", zap.Error(err))
		context.JSON(http.StatusBadRequest, gin.H{"message": "cannot parsed the requested data", "error": true, "data": nil})
		return
	}

	uid, err := user.Save()
	if err != nil {
		utils.Logger.Error("Failed to save user", zap.Error(err))
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not save user", "error": true, "data": nil})
		return
	}

	err = utils.ValidateDetails(user.Name, user.Email, user.Mobile_No, user.Gender, user.Password, user.Confirm_Password)
	if err != nil {
		utils.Logger.Warn("Failed to validate user details", zap.Error(err))
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "error": true, "data": nil})
		return
	}

	userToken, err := utils.GenerateJwtToken(uid)
	if err != nil {
		utils.Logger.Error("Failed to generate user token", zap.Int64("userId", uid), zap.Error(err))
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not generate the token", "error": true, "data": nil})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(uid)
	if err != nil {
		utils.Logger.Error("Failed to generate refresh token", zap.Int64("userId", uid), zap.Error(err))
		context.JSON(http.StatusInternalServerError, gin.H{"message": "culd not generate the token", "error": true, "data": nil})
		return
	}

	err = user.SaveToken(userToken, refreshToken)
	if err != nil {
		utils.Logger.Error("Failed to save token", zap.Int64("userId", uid), zap.Error(err))
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not save the token", "error": true, "data": nil})
		return
	}

	utils.Logger.Info("User signed up successfully", zap.Int64("userId", uid))
	context.JSON(http.StatusCreated, gin.H{"message": "User save successfully", "data": gin.H{"refresh_token": refreshToken, "user_token": userToken}, "error": false})
}

func signIn(context *gin.Context) {
	var login models.Login
	var user models.User

	err := context.ShouldBindJSON(&login)
	if err != nil {
		utils.Logger.Warn("Failed to parse login request", zap.Error(err))
		context.JSON(http.StatusBadRequest, gin.H{"message": "cannot parsed the requested data", "error": true, "data": nil})
		return
	}

	err = login.ValidateCredentials()
	if err != nil {
		utils.Logger.Warn("Authentication failed", zap.Error(err))
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not authenticate user", "error": true, "data": nil})
		return
	}

	userToken, err := utils.GenerateJwtToken(login.ID)
	if err != nil {
		utils.Logger.Error("Failed to generate user token", zap.Int64("userId", login.ID), zap.Error(err))
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not generate the token", "error": true, "data": nil})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(login.ID)
	if err != nil {
		utils.Logger.Error("Failed to generate refresh token", zap.Int64("userId", login.ID), zap.Error(err))
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not generate the token", "error": true, "data": nil})
		return
	}

	err = user.SaveToken(userToken, refreshToken)
	if err != nil {
		utils.Logger.Error("Failed to save token", zap.Int64("userId", login.ID), zap.Error(err))
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not save the token", "error": true, "data": nil})
		return
	}

	utils.Logger.Info("User signed in successfully", zap.Int64("userId", login.ID))
	context.JSON(http.StatusCreated, gin.H{"message": "signIn successfully", "data": gin.H{"refresh_token": refreshToken, "user_token": userToken}, "error": false})
}

func updateUser(c *gin.Context) {
	var req models.UpdateUserRequest

	err := middlewares.CheckTokenPresent(c)
	if err != nil {
		return
	}

	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized: User not authenticated", "error": true, "data": nil})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body", "error": true, "data": nil})
		return
	}

	err = utils.ValidateUser(req.Name, req.Mobile_No)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Data is in invalid format", "error": true, "data": nil})
		return
	}

	err = models.UpdateUserDetails(userId.(int64), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update user", "error": true, "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User details updated successfully", "error": false, "data": nil})
}

func updatePassword(c *gin.Context) {
	var req models.UpdatePasswordRequest

	err := middlewares.CheckTokenPresent(c)
	if err != nil {
		return
	}

	token := c.Request.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")

	// Bind JSON request to struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body", "error": true, "data": nil})
		return
	}

	userIdFromToken, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized: User not authenticated", "error": true, "data": nil})
		return
	}

	err = utils.ValidatePassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to validate password", "error": true, "data": nil})
		return
	}

	user, err := models.GetUserByIdPassChng(userIdFromToken.(int64))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found", "error": true, "data": nil})
		return
	}

	passwordIsValid := utils.CheckPasswordHash(req.OldPassword, user.Password)
	if !passwordIsValid {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Incorrect old password", "error": true, "data": nil})
		return
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to hashed password", "error": true, "data": nil})
		return
	}

	err = models.UpdatePassById(userIdFromToken.(int64), hashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update password", "error": true, "data": nil})
		return
	}

	err = models.DeleteTokenById(userIdFromToken.(int64), token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete tokens", "error": true, "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password Updated Successfully.Please login in again.", "error": false, "data": nil})
}

// signOut function
func signOut(c *gin.Context) {

	tokenString := strings.TrimSpace(c.GetHeader("Authorization"))
	if tokenString == "" {
		utils.Logger.Warn("Token not provided for sign out")
		c.JSON(http.StatusBadRequest, gin.H{"message": "token not found", "error": true, "data": nil})
		return
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	utils.Logger.Info("Recieved signout request")

	err := models.DeleteToken(tokenString)
	if err != nil {
		utils.Logger.Error("Failed to sign out", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to sign out", "error": true, "data": nil})
		return
	}

	utils.Logger.Info("User signed out successfully")
	c.JSON(http.StatusOK, gin.H{"message": "user sign out successfully", "error": false, "data": nil})
}
