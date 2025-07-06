package db

import (
	"database/sql"
	"log"
	"time"

	"github.com/alexleyoung/auto-gcal/internal/types"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func Init() {
	var err error
	db, err = sql.Open("sqlite3", "./agc.db")
	if err != nil {
		log.Fatal(err)
	}

	createUsers := `
    CREATE TABLE IF NOT EXISTS users (
        user_id TEXT NOT NULL PRIMARY KEY,
		email TEXT NOT NULL,
		name TEXT NOT NULL,
		timezone TEXT,
		created_at TEXT DEFAULT (datetime('now'))
    );
	`
	createSessions := `
    CREATE TABLE IF NOT EXISTS sessions (
		id TEXT NOT NULL PRIMARY KEY,
		user_id TEXT NOT NULL,
		access_token TEXT,
		refresh_token TEXT,
		expires_at TEXT,
		FOREIGN KEY(user_id) REFERENCES users(user_id)
    );
    `

	createTables := createUsers + createSessions

	_, err = db.Exec(createTables)
	if err != nil {
		log.Fatal("Failed to create DB:\n" + err.Error())
	}
	log.Println("DB initialized successfully")
}

func GetUser(userID string) (types.User, error) {
	stmt, err := db.Prepare("SELECT * FROM users WHERE user_id = ?")
	if err != nil {
		return types.User{}, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(userID)
	var user types.User
	err = row.Scan(&user.UserID, &user.Email, &user.Name)
	if err != nil {
		return types.User{}, err
	}

	return user, nil
}

func CreateUser(id, email, name, timezone string) (types.User, error) {
	stmt, err := db.Prepare("INSERT INTO users (user_id, email, name, timezone) VALUES (?, ?, ?)")
	if err != nil {
		return types.User{}, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, email, name)
	if err != nil {
		return types.User{}, err
	}

	user := types.User{UserID: id, Email: email, Name: name}
	return user, nil
}

func UpdateUser(id, email, name, timezone string) (types.User, error) {
	stmt, err := db.Prepare("UPDATE users SET email = ?, name = ?, timezone = ? WHERE user_id = ?")
	if err != nil {
		return types.User{}, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(email, name, timezone, id)
	if err != nil {
		return types.User{}, err
	}

	user := types.User{UserID: id, Email: email, Name: name, Timezone: timezone}

	return user, nil
}

func GetSession(sessionID string) (types.Session, error) {
	stmt, err := db.Prepare("SELECT * FROM sessions WHERE id = ?")
	if err != nil {
		return types.Session{}, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(sessionID)
	var session types.Session
	err = row.Scan(&session.ID, &session.UserID, &session.AccessToken, &session.RefreshToken, &session.ExpiresAt)
	if err != nil {
		return types.Session{}, err
	}

	return session, nil
}

func CreateSession(userID, accessToken, refreshToken string, expiresAt time.Time) (types.Session, error) {
	stmt, err := db.Prepare("INSERT INTO sessions (id, user_id, access_token, refresh_token, expires_at) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return types.Session{}, err
	}
	defer stmt.Close()

	// create random session ID
	sessionID := uuid.NewString()

	_, err = stmt.Exec(sessionID, userID, accessToken, refreshToken, expiresAt.Format(time.RFC3339))
	if err != nil {
		return types.Session{}, err
	}
	session := types.Session{ID: sessionID, UserID: userID, AccessToken: accessToken, RefreshToken: refreshToken, ExpiresAt: expiresAt.Format(time.RFC3339)}

	return session, nil
}

func UpdateSessionTokens(sessionID, accessToken string, expiresAt time.Time) error {
	stmt, err := db.Prepare("UPDATE sessions SET access_token = ?, expires_at = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(accessToken, expiresAt.Format(time.RFC3339), sessionID)
	return err
}
