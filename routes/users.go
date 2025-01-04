package routes

import (
	"net/http"
	"strings"
	"task_manager/config"
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
		context.JSON(http.StatusBadRequest, gin.H{"message": "cannot parsed the requested data", "error": true})
		return
	}

	uid, err := user.Save()
	if err != nil {
		utils.Logger.Error("Failed to save user", zap.Error(err))
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not save user", "error": true})
		return
	}

	err = utils.ValidateDetails(user.Name, user.Email, user.Mobile_No, user.Gender, user.Password, user.Confirm_Password)
	if err != nil {
		utils.Logger.Warn("Failed to validate user details", zap.Error(err))
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "error": true})
		return
	}

	userToken, err := utils.GenerateJwtToken(uid)
	if err != nil {
		utils.Logger.Error("Failed to generate user token", zap.Int64("userId", uid), zap.Error(err))
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not generate the token", "error": true})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(uid)
	if err != nil {
		utils.Logger.Error("Failed to generate refresh token", zap.Int64("userId", uid), zap.Error(err))
		context.JSON(http.StatusInternalServerError, gin.H{"message": "culd not generate the token", "error": true})
		return
	}

	err = user.SaveToken(userToken, refreshToken)
	if err != nil {
		utils.Logger.Error("Failed to save token", zap.Int64("userId", uid), zap.Error(err))
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not save the token", "error": true})
		return
	}

	utils.Logger.Info("User signed up successfully", zap.Int64("userId", uid))
	context.JSON(http.StatusCreated, gin.H{"message": "User save successfully", "refresh_token": refreshToken, "user_token": userToken, "error": false})
}

// func getUsers(context *gin.Context) {
// 	users, err := models.GetAllUsers()
// 	if err != nil {
// 		fmt.Println(err)
// 		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch users", "error": true})
// 		return
// 	}

// 	context.JSON(http.StatusOK, users)
// }

// func getLogins(context *gin.Context) {
// 	logins, err := models.GetAllLOGIN()
// 	if err != nil {
// 		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch users"})
// 		return
// 	}

// 	context.JSON(http.StatusOK, logins)
// }

// func getTokens(context *gin.Context) {
// 	tokens, err := models.GetAllToken()
// 	if err != nil {
// 		fmt.Println(err)
// 		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch tokens"})
// 		return
// 	}
// 	context.JSON(http.StatusOK, tokens)
// }

func signIn(context *gin.Context) {
	var login models.Login
	var user models.User

	err := context.ShouldBindJSON(&login)
	if err != nil {
		utils.Logger.Warn("Failed to parse login request", zap.Error(err))
		context.JSON(http.StatusBadRequest, gin.H{"message": "cannot parsed the requested data", "error": true})
		return
	}

	err = login.ValidateCredentials()
	if err != nil {
		utils.Logger.Warn("Authentication failed", zap.Error(err))
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not authenticate user", "error": true})
		return
	}

	userToken, err := utils.GenerateJwtToken(login.ID)
	if err != nil {
		utils.Logger.Error("Failed to generate user token", zap.Int64("userId", login.ID), zap.Error(err))
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not generate the token", "error": true})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(login.ID)
	if err != nil {
		utils.Logger.Error("Failed to generate refresh token", zap.Int64("userId", login.ID), zap.Error(err))
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not generate the token", "error": true})
		return
	}

	err = user.SaveToken(userToken, refreshToken)
	if err != nil {
		utils.Logger.Error("Failed to save token", zap.Int64("userId", login.ID), zap.Error(err))
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not save the token", "error": true})
		return
	}

	utils.Logger.Info("User signed in successfully", zap.Int64("userId", login.ID))
	context.JSON(http.StatusCreated, gin.H{"message": "signIn successfully", "refresh_token": refreshToken, "user_token": userToken, "error": false})
}

func updateUser(c *gin.Context) {
	var updateRequest models.User

	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		utils.Logger.Warn("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body", "error": true})
		return
	}

	if updateRequest.Email == "" {
		utils.Logger.Warn("Email is required")
		c.JSON(http.StatusBadRequest, gin.H{"message": "email is required", "error": true})
		return
	}

	user, err := models.GetUserByEmail(updateRequest.Email)
	if err != nil {
		utils.Logger.Warn("User not found", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"message": "could not get user", "error": true})
		return
	}

	if updateRequest.Name != "" {
		user.Name = updateRequest.Name
		utils.Logger.Info("User name updated successfully")
	}

	if updateRequest.Mobile_No != 0 {
		user.Mobile_No = updateRequest.Mobile_No
		utils.Logger.Info("User mobile number updated successfully")
	}

	if updateRequest.Gender != "" {
		user.Gender = updateRequest.Gender
		utils.Logger.Info("User gender updated successfully")
	}

	if updateRequest.Password != "" {

		if updateRequest.Password != updateRequest.Confirm_Password {
			utils.Logger.Warn("Passwords do not match")
			c.JSON(http.StatusBadRequest, gin.H{"message": "passwords do not match", "error": true})
			return
		}

		hashedPassword, err := utils.HashPassword(updateRequest.Password)
		if err != nil {
			utils.Logger.Error("Failed to hash password", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"message": "could not update user", "error": true})
			return
		}

		if err := config.DB.Model(&config.Login{}).Where("user_id = ?", user.ID).Update("password", hashedPassword).Error; err != nil {
			utils.Logger.Error("Failed to update password", zap.Int64("user_id", user.ID), zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update password", "error": true})
			return
		}

		utils.Logger.Info("User password updated successfully", zap.Int64("user_id", user.ID))
	}

	err = updateRequest.UpdateUserTable()
	if err != nil {
		utils.Logger.Error("Failed to update user", zap.Int64("user_id", user.ID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not update user", "error": true})
		return
	}

	utils.Logger.Info("User updated successfully", zap.Int64("user_id", user.ID))
	c.JSON(http.StatusOK, gin.H{"message": "user updated successfully", "error": false})
}

// signOut function
func signOut(c *gin.Context) {

	tokenString := strings.TrimSpace(c.GetHeader("Authorization"))
	if tokenString == "" {
		utils.Logger.Warn("Token not provided for sign out")
		c.JSON(http.StatusBadRequest, gin.H{"message": "token not found", "error": true})
		return
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	utils.Logger.Info("Recieved signout request")

	err := models.DeleteToken(tokenString)
	if err != nil {
		utils.Logger.Error("Failed to sign out", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to sign out", "error": true})
		return
	}

	utils.Logger.Info("User signed out successfully")
	c.JSON(http.StatusOK, gin.H{"message": "user sign out successfully", "error": false})
}
