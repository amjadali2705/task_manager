package routes

import (
	"fmt"
	"net/http"
	"task_manager/config"
	"task_manager/models"
	"task_manager/utils"

	"github.com/gin-gonic/gin"
)

func signUp(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "cannot parsed the requested data", "error": true})
		return
	}

	uid, err := user.Save()
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not save user", "error": true})
		return
	}

	err = utils.ValidateDetails(user.Name, user.Email, user.Mobile_No, user.Gender, user.Password, user.Confirm_Password)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "error": true})
		return
	}

	jwtToken, err := utils.GenerateJwtToken(uid)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not generate the token", "error": true})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(uid)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "culd not generate the token", "error": true})
		return
	}

	err = user.SaveToken(jwtToken, refreshToken)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not save the token", "error": true})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "User save successfully", "refresh_token": refreshToken, "user_token": jwtToken, "error": false})
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

	err := context.ShouldBindJSON(&login)
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "cannot parsed the requested data", "error": true})
		return
	}

	err = login.ValidateCredentials()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not authenticate user", "error": true})
		return
	}

	jwtToken, err := utils.GenerateJwtToken(login.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not generate the token", "error": true})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(login.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not generate the token", "error": true})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "signIn successfully", "refresh_token": refreshToken, "user_token": jwtToken, "error": false})
}

func updateUser(c *gin.Context) {
	var updateRequest models.User

	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "could not parse request", "error": true})
		return
	}

	if updateRequest.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "email is required", "error": true})
		return
	}

	user, err := models.GetUserByEmail(updateRequest.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "could not get user", "error": true})
		return
	}

	if updateRequest.Name != "" {
		user.Name = updateRequest.Name
	}

	if updateRequest.Mobile_No != 0 {
		user.Mobile_No = updateRequest.Mobile_No
	}

	if updateRequest.Gender != "" {
		user.Gender = updateRequest.Gender
	}

	if updateRequest.Password != "" {

		if updateRequest.Password != updateRequest.Confirm_Password {
			c.JSON(http.StatusBadRequest, gin.H{"message": "passwords do not match", "error": true})
			return
		}

		hashedPassword, err := utils.HashPassword(updateRequest.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "could not update user", "error": true})
			return
		}

		if err := config.DB.Model(models.Login{}).Where("user_id = ?", user.ID).Update("password", hashedPassword).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update user", "error": true})
			return
		}
	}

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not update user", "error": true})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user updated successfully", "error": false})
}
