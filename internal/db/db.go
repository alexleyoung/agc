package db

import (
	"database/sql"
	"log"

	"github.com/alexleyoung/auto-gcal/internal/types"
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
        user_id TEXT NOT NULL PRIMARY KEY,
		email TEXT NOT NULL
    );
    CREATE TABLE IF NOT EXISTS auth (
        user_id INTEGER NOT NULL PRIMARY KEY,
		token TEXT,
		FOREIGN KEY(user_id) REFERENCES users(id)
    );
    `

	_, err = db.Exec(createTables)
	if err != nil {
		log.Fatal("Failed to create DB:\n" + err.Error())
	}
	log.Println("DB initialized successfully")
}

func CreateUser(id, email string) error {
	stmt, err := db.Prepare("INSERT INTO users (user_id, email) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, email)
	if err != nil {
		return err
	}

	return nil
}

func SaveToken(userInfo types.UserInfo, token *oauth2.Token) error {
	return nil
}
