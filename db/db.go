package db

import "database/sql"
// import _ "github.com/mattn/go-sqlite3"
import _ "github.com/glebarez/sqlite"

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite", "api.db")

	if err != nil {
		panic("error connecting to database")
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	createTables()
}

func createTables() {
	createUserDetailsTable := `CREATE
		TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			mob_no INTEGER NOT NULL
		);`

	_, err := DB.Exec(createUserDetailsTable)
	if err != nil {
		panic("error creating users table")
	}

	createUserPasswordTable := `CREATE TABLE IF NOT EXISTS user_passwords (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		password TEXT NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`

	_, err = DB.Exec(createUserPasswordTable)
	if err != nil {
		panic("error creating user_passwords table")
	}
}
