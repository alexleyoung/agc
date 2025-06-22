package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/oauth2"
)

var db *sql.DB

func Init() {
	db, err := sql.Open("sqlite3", "./agc.db")
	if err != nil {
		log.Fatal(err)
	}

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
	stmt, err := db.Prepare("INSERT INTO users (email) VALUES (?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(email)
	if err != nil {
		return err
	}

	return nil
}

func SaveToken(userID string, token *oauth2.Token) error {
	return nil
}
