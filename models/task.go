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

func GetAllTasks() ([]Task, error) {
	query := `SELECT * FROM tasks`
	rows, err := config.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.IsCompleted, &task.UserID)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func GetTaskById(id int64) (*Task, error) {
	query := `SELECT * FROM tasks WHERE id = ?`
	row := config.DB.QueryRow(query, id)

	var task Task

	err := row.Scan(&task.ID, &task.Title, &task.Description, &task.IsCompleted, &task.UserID)
	if err != nil {	
		return nil, err
	}

	return &task, nil
}

func (t *Task) Update() error {
	query := `UPDATE tasks SET title = ?, description = ?, isCompleted = ? WHERE id = ?`

	stmt, err := config.DB.Prepare(query)
	if err != nil {
		return err
	}		
	defer stmt.Close()

	_, err = stmt.Exec(t.Title, t.Description, t.IsCompleted, t.ID)
	return err
}

