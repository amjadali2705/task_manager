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
	"time"
)

type Task struct {
	ID          int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string `gorm:"not null" json:"title" binding:"required"`
	Description string `gorm:"not null" json:"description" binding:"required"`
	IsCompleted string `gorm:"not null" json:"isCompleted" binding:"required"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	UserID      int64 `gorm:"not null" json:"userId"`
}

func (t *Task) Save() error {
	result := config.DB.Create(t)
	return result.Error
}

func GetAllTasks() ([]Task, error) {
	var tasks []Task
	result := config.DB.Find(&tasks)
	return tasks, result.Error
}

func GetTaskByID(id int64) (*Task, error) {
	var task Task
	result := config.DB.First(&task, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &task, nil
}

func (t *Task) Update() error {
	result := config.DB.Model(&Task{}).Where("id = ?", t.ID).Updates(Task{Title: t.Title, Description: t.Description, IsCompleted: t.IsCompleted, UpdatedAt: time.Now()})
	return result.Error
}

func (t *Task) Delete() error {
	result := config.DB.Delete(t)
	return result.Error
}

func GetTasksWithFilters(userId int64, sortOrder, isCompleted string, limit, offset int) ([]Task, int64, error) {
	var tasks []Task
	var totalTasks int64

	// Start building the query
	query := config.DB.Order("created_at " + sortOrder).Where("user_id = ?", userId)

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
	return tasks, totalTasks, result.Error
}

