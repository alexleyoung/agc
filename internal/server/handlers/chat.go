package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/alexleyoung/agc/internal/ai"
	"github.com/alexleyoung/agc/internal/types"
)

func setupAI(mux *http.ServeMux) {
	mux.Handle("POST /chat", AuthMiddleware(http.HandlerFunc(chat)))
	mux.HandleFunc("GET /chat", test)
}

func test(w http.ResponseWriter, _ *http.Request) {
	tmpl, err := template.ParseFiles("internal/templates/home.html")
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func chat(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(sessionKey).(types.Session)
	if !ok {
		http.Error(w, "Session missing in context", http.StatusInternalServerError)
		return
	}

	var body types.ChatRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Printf("Error parsing body: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := ai.Chat(r.Context(), session, body.Model, body.History, body.Prompt)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
	}
}
