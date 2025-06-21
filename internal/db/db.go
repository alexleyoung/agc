package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func Init() {
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTables := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		email TEXT
    );
    CREATE TABLE IF NOT EXISTS auth (
        user_id INTEGER NOT NULL PRIMARY KEY FOREIGN KEY REFERENCES users(id),
		token TEXT
    );
    `
	_, err = db.Exec(createTables)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Table 'users' created successfully")
}
