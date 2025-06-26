package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/alexleyoung/auto-gcal/internal/ai"
	"google.golang.org/genai"
)

type chatRequestBody struct {
	Prompt  string           `json:"prompt"`
	Model   string           `json:"model,omitempty"`
	History []*genai.Content `json:"history,omitempty"`
}

func setupAI(mux *http.ServeMux) {
	mux.HandleFunc("GET /chat", http.HandlerFunc(chat))
}

func chat(w http.ResponseWriter, r *http.Request) {
	var body chatRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Printf("Error parsing body: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	session, err := r.Cookie("agc_session")
	if err != nil {
		log.Print("Error getting session cookie:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := session.Value

	// SCUFFED FOR SIMPLE TESTING - SETUP ARGS
	prompt := r.URL.Query().Get("prompt")
	if prompt == "" {
		http.Error(w, "No prompt submitted.", http.StatusBadRequest)
		return
	}
	model := body.Model
	if model == "" {
		model = "gemini-2.0-flash"
	}
	history := body.History
	if history == nil {
		history = make([]*genai.Content, 0)
	}

	res, err := ai.Chat(r.Context(), userID, model, history, prompt)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.Write([]byte(res.Text()))
	}
}
