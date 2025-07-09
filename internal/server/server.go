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
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading environment: %v", err)
	}
	log.Print("Successfully loaded environment")
	PORT := os.Getenv("PORT")

	db.Init()
	log.Print("Initialized database")
	auth.Init()
	log.Print("Initialized OAuth2 client")

	mux := http.NewServeMux()
	handlers.Init(mux)
	log.Print("Initialized handlers")

	log.Println("Server starting on port :" + PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, mux))
}
