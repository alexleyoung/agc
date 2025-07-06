package server

import (
	"log"
	"net/http"
	"os"

	"github.com/alexleyoung/auto-gcal/internal/auth"
	"github.com/alexleyoung/auto-gcal/internal/db"
	"github.com/alexleyoung/auto-gcal/internal/server/handlers"
	"github.com/joho/godotenv"
)

func Run() {
	log.Print("Loading environment variables...")
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading environment: %v", err)
	}
	PORT := os.Getenv("PORT")

	log.Print("Initializing database...")
	db.Init()
	log.Print("Initializing OAuth2 client...")
	auth.Init()

	log.Print("Initializing server...")
	mux := http.NewServeMux()
	handlers.Init(mux)

	log.Println("Server starting on port :" + PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, mux))
}
