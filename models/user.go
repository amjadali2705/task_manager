package models

import "task_manager/db"

type User struct {
	ID         int
	Email      string `json:"email"`
	FName      string `json:"first_name"`
	LName      string `json:"last_name"`
	MobNo      int    `json:"mob_no"`
	// ProfilePic string `binding:"required"`
}

type UserPassword struct {
	ID       int
	UserID   int
	Password string `binding:"required"`
	Emaill   string `binding:"required"`
}

func (u *User) Save() error {
	query := `INSERT INTO users (first_name, last_name, email, mob_no) VALUES (?, ?, ?, ?)`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return nil

}
