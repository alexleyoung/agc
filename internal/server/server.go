package server

import (
	"log"
	"net/http"
	"os"

	"github.com/alexleyoung/auto-gcal/internal/server/handlers"
	"github.com/joho/godotenv"
)

func Run() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading environment: %v", err)
	}
	PORT := os.Getenv("PORT")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /chat", handlers.Chat)

	log.Println("Server starting on port :" + PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, mux))

}
