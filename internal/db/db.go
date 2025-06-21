package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func Init() {
	db, err := sql.Open("sqlite3", "./agc.db")
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
        user_id INTEGER NOT NULL PRIMARY KEY,
		token TEXT,
		FOREIGN KEY(user_id) REFERENCES users(id)
    );
    `

	_, err = db.Exec(createTables)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("DB initialized successfully")
}

func CreateUser(email string) error {
	return nil
}
