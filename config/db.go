// package config

// import (
// 	"database/sql"
// 	"fmt"

// 	_ "github.com/glebarez/sqlite"
// )

// var DB *sql.DB

// func InitDB() {
// 	var err error
// 	DB, err = sql.Open("sqlite", "task.db")

// 	if err != nil {
// 		fmt.Println(err)
// 		panic("Could not connect to database.")
// 	}

// 	DB.SetMaxOpenConns(10)
// 	DB.SetMaxIdleConns(5)

// 	createTables()
// }

// func createTables() {
// 	createUserTable := ` CREATE TABLE IF NOT EXISTS users(
// 		u_id INTEGER PRIMARY KEY AUTOINCREMENT,
// 		name TEXT NOT NULL,
// 		mobile_no INTEGER NOT NULL,
// 		gender TEXT NOT NULL,
// 		email TEXT NOT NULL UNIQUE
// 	) `

// 	_, err := DB.Exec(createUserTable)
// 	if err != nil {
// 		panic("Could not create users table")
// 	}

// 	createLoginTable := ` CREATE TABLE IF NOT EXISTS login(
// 		id INTEGER PRIMARY KEY AUTOINCREMENT,
// 		email TEXT NOT NULL UNIQUE,
// 		password TEXT NOT NULL,
// 		user_id INTEGER,
// 		FOREIGN KEY (user_id) REFERENCES users(u_id)
// 	) `

// 	_, err = DB.Exec(createLoginTable)
// 	if err != nil {
// 		panic("Could not create login table")
// 	}

// 	createTokenTable := ` CREATE TABLE IF NOT EXISTS token(
// 		id INTEGER PRIMARY KEY AUTOINCREMENT,
// 		refresh_token TEXT NOT NULL UNIQUE,
// 		jwt_token TEXT NOT NULL UNIQUE,
// 		timestamp TIMESTAMP NOT NULL
// 	) `

// 	_, err = DB.Exec(createTokenTable)
// 	if err != nil {
// 		panic("Could not create token table")
// 	}

// 	createTaskTable := ` CREATE TABLE IF NOT EXISTS tasks(
// 		id INTEGER PRIMARY KEY AUTOINCREMENT,
// 		title TEXT NOT NULL,
// 		description TEXT NOT NULL,
// 		isCompleted VARCHAR(5) NOT NULL,
// 		user_id INTEGER,
// 		FOREIGN KEY (user_id) REFERENCES users(u_id)
// 	) `

// 	_, err = DB.Exec(createTaskTable)
// 	if err != nil {
// 		panic("Could not create tasks table")
// 	}

// }

package config

import (
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("task.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	createTables()
}

func createTables() {
	err := DB.AutoMigrate(&User{}, &Login{}, &Token{}, &Task{})
	if err != nil {
		log.Fatalf("could not migrate tables: %v", err)
	}

}

type User struct {
	ID       int64  `gorm:"primaryKey;autoIncrement"`
	Name     string `gorm:"not null"`
	MobileNo int64  `gorm:"not null"`
	Gender   string `gorm:"not null"`
	Email    string `gorm:"not null;unique"`
}

type Login struct {
	ID       int64  `gorm:"primaryKey;autoIncrement"`
	Email    string `gorm:"not null;unique"`
	Password string `gorm:"not null"`
	UserID   int64
	User     User `gorm:"foreignKey:UserID"`
}

type Token struct {
	ID           int64     `gorm:"primaryKey;autoIncrement"`
	RefreshToken string    `gorm:"not null;unique"`
	UserToken    string    `gorm:"not null;unique"`
	Timestamp    time.Time `gorm:"not null"`
}

type Task struct {
	ID          int64  `gorm:"primaryKey;autoIncrement"`
	Title       string `gorm:"not null"`
	Description string `gorm:"not null"`
	IsCompleted string `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	UserID      int64
	User        User `gorm:"foreignKey:UserID"`
}
