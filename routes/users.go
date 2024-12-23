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
		context.JSON(http.StatusBadRequest, gin.H{"message": "cannot parsed the requested data"})
		return
	}

	err = user.Save()
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not save user"})
		return
	}

	jwtToken, err := utils.GenerateJwtToken(user.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "culd not generate the token"})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "culd not generate the token"})
		return
	}

	err = user.SaveToken(jwtToken, refreshToken)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "culd not save the token"})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "User save successfully", "refresh": refreshToken, "jwt": jwtToken})
}

func getUsers(context *gin.Context) {
	users, err := models.GetAllUsers()
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch users"})
		return
	}

	context.JSON(http.StatusOK, users)
}

func getLogins(context *gin.Context) {
	logins, err := models.GetAllLOGIN()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch users"})
		return
	}

	context.JSON(http.StatusOK, logins)
}

func getTokens(context *gin.Context) {
	tokens, err := models.GetAllToken()
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch tokens"})
		return
	}
	context.JSON(http.StatusOK, tokens)
}
