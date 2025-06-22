package auth

import (
	"net/http"

	"github.com/alexleyoung/auto-gcal/internal/auth"
	"golang.org/x/oauth2"
)

func GetAuthURL(w http.ResponseWriter, r *http.Request) {
	config := auth.GetConfig()
	url := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusFound)
}

func Callback(w http.ResponseWriter, r *http.Request) {
	config := auth.GetConfig()
	tok, err := config.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(tok.AccessToken))
}
