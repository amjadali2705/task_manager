package models

import "task_manager/config"

type Task struct {
	ID          int64
	Title       string `binding:"required" json:"title"`
	Description string `binding:"required" json:"description"`
	IsCompleted string `binding:"required" json:"isCompleted"`
	UserID      int64
}

func (t *Task) Save() error {
	query := `INSERT INTO tasks (title, description, isCompleted, user_id) VALUES (?, ?, ?, ?)`

	stmt, err := config.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(t.Title, t.Description, t.IsCompleted, t.UserID)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	t.ID = id
	return err
}
