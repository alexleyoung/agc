package db

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"

	"github.com/alexleyoung/auto-gcal/internal/types"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/oauth2"
)

var db *sql.DB

func Init() {
	var err error
	db, err = sql.Open("sqlite3", "./agc.db")
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
		FOREIGN KEY(user_id) REFERENCES users(user_id)
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

func GetUserToken(userID string) (string, error)

func SaveToken(userInfo types.UserInfo, token *oauth2.Token) error {
	stmt, err := db.Prepare("INSERT INTO auth (user_id, token) VALUES (?, ?);")
	if err != nil {
		return err
	}
	defer stmt.Close()

	encToken, err := encryptToken(token)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(userInfo.Sub, encToken)
	if err != nil {
		return err
	}
	return nil
}

func encryptToken(token *oauth2.Token) (string, error) {
	data, err := json.Marshal(token)
	if err != nil {
		return "", nil
	}

	ciphertext, err := encrypt(data, []byte(os.Getenv("OAUTH_TOKEN_CIPHER_KEY")))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func DecryptToken(encryptedToken string) (*oauth2.Token, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedToken)
	if err != nil {
		return nil, err
	}

	tokBlob, err := decrypt(ciphertext, []byte(os.Getenv("OAUTH_TOKEN_CIPHER_KEY")))
	if err != nil {
		return nil, err
	}

	var tok oauth2.Token
	err = json.Unmarshal(tokBlob, &tok)
	if err != nil {
		return nil, err
	}

	return &tok, nil
}

func encrypt(plaintext, key []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("Key must be 32 bytes for AES-256")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

func decrypt(ciphertext, key []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("Key must be 32 bytes for AES-256")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce := ciphertext[:nonceSize]
	ciphertextData := ciphertext[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertextData, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
