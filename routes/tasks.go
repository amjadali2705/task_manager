package routes

import (
	"net/http"
	"strconv"
	"task_manager/middlewares"
	"task_manager/models"
	"task_manager/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func createTask(context *gin.Context) {

	err := middlewares.CheckTokenPresent(context)
	if err != nil {
		return
	}

	var task models.Task
	err = context.ShouldBindJSON(&task)
	if err != nil {
		utils.Logger.Error("Failed to parse request", zap.Error(err))
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse request", "error": true, "data": nil})
		return
	}

	userId := context.GetInt64("userId")
	task.UserID = userId
	utils.Logger.Info("Recieved task creation request", zap.Int64("userId", userId))

	err = task.Save()
	if err != nil {
		utils.Logger.Error("Failed to save task", zap.Error(err))
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not create task", "error": true, "data": nil})
		return
	}

	utils.Logger.Info("Task created successfully", zap.Int64("taskId", task.ID), zap.Int64("userId", userId))
	context.JSON(http.StatusCreated, gin.H{"message": "task created", "data": gin.H{"taskId": task.ID}, "error": false})
}

// func getTasks(context *gin.Context) {
// 	tasks, err := models.GetAllTasks()
// 	if err != nil {
// 		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not get tasks", "error": true})
// 		return
// 	}

// 	context.JSON(http.StatusOK, gin.H{"tasks": tasks, "error": false})
// }

func getTask(context *gin.Context) {

	err := middlewares.CheckTokenPresent(context)
	if err != nil {
		return
	}

	userId, exists := context.Get("userId")
	if !exists {
		utils.Logger.Error("User ID not found in context", zap.String("context", "userId"))
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Not Authorized", "error": true, "data": nil})
		return
	}

	taskId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		utils.Logger.Error("Failed to parse task ID", zap.String("param", context.Param("id")), zap.Error(err))
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse to int", "error": true, "data": nil})
		return
	}

	utils.Logger.Info("Fetching task", zap.Int64("taskId", taskId), zap.Int64("userId", userId.(int64)))

	task, err := models.GetTaskByID(taskId, userId.(int64))
	if err != nil {
		utils.Logger.Error("Failed to get task or access denied", zap.Int64("taskId", taskId), zap.Int64("userId", userId.(int64)))
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not get task", "error": true, "data": nil})
		return
	}

	utils.Logger.Info("Task fetched successfully", zap.Int64("taskId", taskId), zap.Int64("userId", userId.(int64)))
	context.JSON(http.StatusOK, gin.H{"message": "Task fetched successfully", "data": task, "error": false})
}

func updateTask(context *gin.Context) {
	err := middlewares.CheckTokenPresent(context)
	if err != nil {
		return
	}

	taskId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		utils.Logger.Error("Failed to parse task ID", zap.String("param", context.Param("id")), zap.Error(err))
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse to int", "error": true, "data": nil})
		return
	}

	userID := context.GetInt64("userId")

	task, err := models.GetTaskByID(taskId, userID)
	if err != nil {
		utils.Logger.Error("Failed to fetch task", zap.Int64("taskId", taskId), zap.Int64("userId", userID), zap.Error(err))
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not get task", "error": true, "data": nil})
		return
	}

	if task.UserID != userID {
		utils.Logger.Warn("User not authorized to update task", zap.Int64("taskId", taskId), zap.Int64("userId", userID))
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Not Authorized to update", "error": true, "data": nil})
		return
	}

	var updatedTask models.Task

	err = context.ShouldBindJSON(&updatedTask)
	if err != nil {
		utils.Logger.Error("Failed to bind json", zap.Error(err))
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse request", "error": true, "data": nil})
		return
	}

	updatedTask.ID = taskId

	err = updatedTask.Update()
	if err != nil {
		utils.Logger.Error("Failed to update task", zap.Int64("taskId", taskId), zap.Error(err))
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not update task", "error": true, "data": nil})
		return
	}

	utils.Logger.Info("Task updated successfully", zap.Int64("taskId", taskId), zap.Int64("userId", userID))
	context.JSON(http.StatusOK, gin.H{"message": "task updated successfully", "error": false, "data": nil})
}

func deleteTask(context *gin.Context) {
	err := middlewares.CheckTokenPresent(context)
	if err != nil {
		return
	}

	taskId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		utils.Logger.Error("Failed to parse task ID", zap.String("param", context.Param("id")), zap.Error(err))
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse to int", "error": true, "data": nil})
		return
	}

	userID := context.GetInt64("userId")

	task, err := models.GetTaskByID(taskId, userID)
	if err != nil {
		utils.Logger.Error("Failed to fetch task for deletion", zap.Int64("taskId", taskId), zap.Int64("userId", userID), zap.Error(err))
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not get task", "error": true, "data": nil})
		return
	}

	if task.UserID != userID {
		utils.Logger.Warn("User not authorized to delete task", zap.Int64("taskId", taskId), zap.Int64("userId", userID))
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Not Authorized to delete", "error": true, "data": nil})
		return
	}

	err = task.Delete()
	if err != nil {
		utils.Logger.Error("Failed to delete task", zap.Int64("taskId", taskId), zap.Error(err))
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not delete task", "error": true, "data": nil})
		return
	}

	utils.Logger.Info("Task deleted successfully", zap.Int64("taskId", taskId), zap.Int64("userId", userID))
	context.JSON(http.StatusOK, gin.H{"message": "task deleted successfully", "error": false, "data": nil})
}

func getTasksByQuery(context *gin.Context) {

	err := middlewares.CheckTokenPresent(context)
	if err != nil {
		return
	}

	userId, exists := context.Get("userId")
	if !exists {
		utils.Logger.Warn("Unauthorized access attempt in getTasksByQuery")
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Not Authorized", "error": true, "data": nil})
		return
	}

	// Retrieve query parameters
	sortOrder := context.DefaultQuery("sort", "asc") // Default sort order: ascending
	isCompleted := context.Query("isCompleted")      // Optional filter

	// Pagination parameters
	page, err := strconv.Atoi(context.DefaultQuery("page", "1")) // Default page: 1
	if err != nil || page < 1 {
		utils.Logger.Warn("Invalid page parameter", zap.String("page", context.DefaultQuery("page", "1")), zap.Error(err))
		page = 1
	}

	limit, err := strconv.Atoi(context.DefaultQuery("limit", "10")) // Default limit: 10
	if err != nil || limit < 1 {
		utils.Logger.Warn("Invalid limit parameter", zap.String("limit", context.DefaultQuery("limit", "5")), zap.Error(err))
		limit = 5
	}

	offset := (page - 1) * limit

	// Fetch tasks with filters, sorting, and pagination
	tasks, totalTasks, err := models.GetTasksWithFilters(userId.(int64), sortOrder, isCompleted, limit, offset)
	if err != nil {
		utils.Logger.Error("Failed to get tasks", zap.Int64("userId", userId.(int64)), zap.Error(err))
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not get tasks", "error": true, "data": nil})
		return
	}

	// Calculate total pages
	totalPages := (totalTasks + int64(limit) - 1) / int64(limit)

	// Respond with tasks and pagination metadata
	utils.Logger.Info("Tasks fetched successfully", zap.Int64("userId", userId.(int64)), zap.String("sortOrder", sortOrder), zap.String("isCompleted", isCompleted), zap.Int("page", page), zap.Int("limit", limit), zap.Int("totalPages", int(totalPages)))
	context.JSON(http.StatusOK, gin.H{
		"message": "Tasks fetched successfully",
		"data": gin.H{
		"tasks":       tasks,
		"totalPages":  totalPages,
		"currentPage": page},
		"error":       false,
	})
}
