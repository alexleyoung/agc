package server

import (
	"log"
	"net/http"
	"os"

	"github.com/alexleyoung/agc/internal/auth"
	"github.com/alexleyoung/agc/internal/db"
	"github.com/alexleyoung/agc/internal/server/handlers"
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
	log.Print("Successfully initialized database")
	auth.Init()
	log.Print("Successfully initialized OAuth2 client")

	mux := http.NewServeMux()
	handlers.Init(mux)
	log.Print("Successfully initialized handlers")

	log.Println("Server starting on port :" + PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, mux))
}
