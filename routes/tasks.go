package routes

import (
	"net/http"
	"task_manager/models"

	"github.com/gin-gonic/gin"
)

func createTask(context *gin.Context) {
	var task models.Task
	err := context.ShouldBindJSON(&task)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse request", "error": true})
		return
	}

	userId := context.GetInt64("userId")
	task.UserID = userId

	err = task.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not create task", "error": true})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "task created", "tasks": task, "error": false})
}

func getTasks(context *gin.Context) {
	tasks, err := models.GetAllTasks()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not get tasks", "error": true})
		return
	}

	context.JSON(http.StatusOK, gin.H{"tasks": tasks, "error": false})
}
