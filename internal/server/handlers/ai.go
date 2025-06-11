package handlers

import (
	"log"
	"net/http"

	"github.com/alexleyoung/auto-gcal/internal/ai"
)

func Chat(w http.ResponseWriter, r *http.Request) {
	prompt := r.URL.Query().Get("prompt")
	if prompt == "" {
		http.Error(w, "No prompt submitted.", http.StatusBadRequest)
		return
	}

	res, err := ai.Chat(r.Context(), "gemini-2.0-flash", prompt)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.Write([]byte(res.Text()))
	}
}
