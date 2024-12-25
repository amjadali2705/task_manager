package routes

import (
	"fmt"
	"net/http"
	"task_manager/models"
	"task_manager/utils"

	"github.com/gin-gonic/gin"
)

func signup(context *gin.Context) {
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

	err = utils.ValidateDetails(user.Name, user.Email, user.Mobile_No, user.Gender, user.Password)
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

	context.JSON(http.StatusCreated, gin.H{"message": "User save successfully", "refresh_token": refreshToken, "access_token": jwtToken, "error": false})
}

func getUsers(context *gin.Context) {
	users, err := models.GetAllUsers()
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch users", "error": true})
		return
	}

	context.JSON(http.StatusOK, users)
}

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

func login(context *gin.Context) {
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

	context.JSON(http.StatusOK, gin.H{"message": "login successfully", "refresh_token": refreshToken, "access_token": jwtToken, "error": false})
}
