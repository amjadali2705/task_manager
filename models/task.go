// package models

// import "task_manager/config"

// type Task struct {
// 	ID          int64
// 	Title       string `binding:"required" json:"title"`
// 	Description string `binding:"required" json:"description"`
// 	IsCompleted string `binding:"required" json:"isCompleted"`
// 	UserID      int64
// }

// func (t *Task) Save() error {
// 	query := `INSERT INTO tasks (title, description, isCompleted, user_id) VALUES (?, ?, ?, ?)`

// 	stmt, err := config.DB.Prepare(query)
// 	if err != nil {
// 		return err
// 	}
// 	defer stmt.Close()

// 	result, err := stmt.Exec(t.Title, t.Description, t.IsCompleted, t.UserID)
// 	if err != nil {
// 		return err
// 	}

// 	id, err := result.LastInsertId()
// 	t.ID = id
// 	return err
// }

// func GetAllTasks() ([]Task, error) {
// 	query := `SELECT * FROM tasks`
// 	rows, err := config.DB.Query(query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var tasks []Task
// 	for rows.Next() {
// 		var task Task
// 		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.IsCompleted, &task.UserID)
// 		if err != nil {
// 			return nil, err
// 		}
// 		tasks = append(tasks, task)
// 	}
// 	return tasks, nil
// }

// func GetTaskById(id int64) (*Task, error) {
// 	query := `SELECT * FROM tasks WHERE id = ?`
// 	row := config.DB.QueryRow(query, id)

// 	var task Task

// 	err := row.Scan(&task.ID, &task.Title, &task.Description, &task.IsCompleted, &task.UserID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &task, nil
// }

// func (t *Task) Update() error {
// 	query := `UPDATE tasks SET title = ?, description = ?, isCompleted = ? WHERE id = ?`

// 	stmt, err := config.DB.Prepare(query)
// 	if err != nil {
// 		return err
// 	}
// 	defer stmt.Close()

// 	_, err = stmt.Exec(t.Title, t.Description, t.IsCompleted, t.ID)
// 	return err
// }

// func (t *Task) Delete() error {
// 	query := `DELETE FROM tasks WHERE id = ?`

// 	stmt, err := config.DB.Prepare(query)
// 	if err != nil {
// 		return err
// 	}
// 	defer stmt.Close()

// 	_, err = stmt.Exec(t.ID)
// 	return err
// }

package models

import (
	"task_manager/config"
	"task_manager/utils"
	"time"

	"go.uber.org/zap"
)

type Task struct {
	ID          int64  `json:"id"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	IsCompleted string `json:"isCompleted" binding:"required"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	UserID      int64 `json:"userId"`
}

func (t *Task) Save() error {
	result := config.DB.Create(&t)
	if result.Error != nil {
		utils.Logger.Error("Failed to save task", zap.Error(result.Error), zap.Int64("userId", t.ID))
		return result.Error
	}

	utils.Logger.Info("Task saved successfully", zap.Int64("taskId", t.ID))
	return nil
}

// func GetAllTasks() ([]Task, error) {
// 	var tasks []Task
// 	result := config.DB.Find(&tasks)
// 	return tasks, result.Error
// }

func GetTaskByID(id, userId int64) (*Task, error) {
	var task Task
	result := config.DB.Where("id = ? AND user_id = ?", id, userId).First(&task)
	if result.Error != nil {
		utils.Logger.Error("Failed to fetch task by id", zap.Error(result.Error), zap.Int64("taskId", id))
		return &Task{}, result.Error
	}

	utils.Logger.Info("Task fetched by id successfully", zap.Int64("taskId", id), zap.Int64("userId", userId))
	return &task, nil
}

func (t *Task) Update() error {
	result := config.DB.Model(&Task{}).Where("id = ?", t.ID).Updates(Task{Title: t.Title, Description: t.Description, IsCompleted: t.IsCompleted, UpdatedAt: time.Now()})
	if result.Error != nil {
		utils.Logger.Error("Failed to update task", zap.Error(result.Error), zap.Int64("taskId", t.ID))
		return result.Error
	}

	utils.Logger.Info("Task updated successfully", zap.Int64("taskId", t.ID))
	return result.Error
}

func (t *Task) Delete() error {
	result := config.DB.Delete(t)
	if result.Error != nil {
		utils.Logger.Error("Failed to delete task", zap.Error(result.Error), zap.Int64("taskId", t.ID))
		return result.Error
	}

	utils.Logger.Info("Task deleted successfully", zap.Int64("taskId", t.ID))
	return nil
}

func GetTasksWithFilters(userId int64, sortOrder, isCompleted string, limit, offset int) ([]Task, int64, error) {
	var tasks []Task
	var totalTasks int64

	// Start building the query
	query := config.DB.Order("created_at "+sortOrder).Where("user_id = ?", userId)

	// Apply the isCompleted filter if provided
	if isCompleted != "" {
		query = query.Where("is_completed = ?", isCompleted)
	}

	// Count the total number of tasks (without limit/offset)
	query.Model(&Task{}).Count(&totalTasks)

	// Apply pagination
	query = query.Limit(limit).Offset(offset)

	// Execute the query
	result := query.Find(&tasks)
	if result.Error != nil {
		utils.Logger.Error("Failed to fetch tasks", zap.Error(result.Error), zap.Int64("userId", userId), zap.String("sortOrder", sortOrder), zap.String("isCompleted", isCompleted), zap.Int("limit", limit), zap.Int("offset", offset))
		return nil, 0, result.Error
	}

	utils.Logger.Info("Tasks fetched successfully", zap.Int64("userId", userId), zap.String("sortOrder", sortOrder), zap.String("isCompleted", isCompleted), zap.Int("tasksCount", len(tasks)), zap.Int("totalTasks", int(totalTasks)))
	return tasks, totalTasks, nil
}
