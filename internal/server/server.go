package server

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func Run() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading environment: %v", err)
	}
	PORT := os.Getenv("PORT")

	mux := http.NewServeMux()

	// mux.HandleFunc("GET /", pages.Home)
	// mux.HandleFunc("GET /blab/{slug}", pages.Post)
	// mux.HandleFunc("GET /yap", pages.Create)
	// mux.HandleFunc("GET /edit", pages.Edit)
	//
	// mux.HandleFunc("POST /blab", handlers.PostPost)
	// mux.HandleFunc("PATCH /blab", handlers.PatchPost)

	log.Println("Server starting on port :" + PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, mux))

}
