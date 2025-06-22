package handlers

import (
	"net/http"

	"github.com/alexleyoung/auto-gcal/internal/auth"
)

func setupAuth(mux *http.ServeMux) {
	mux.HandleFunc("GET /auth", getAuthURL)
	mux.HandleFunc("GET /auth/callback", oauthCallback)
}

func getAuthURL(w http.ResponseWriter, r *http.Request) {
	url := auth.GetAuthURL()
	http.Redirect(w, r, url, http.StatusFound)
}

func oauthCallback(w http.ResponseWriter, r *http.Request) {
	info, err := auth.Authenticate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Authenticated as user: " + info.Email))
}
