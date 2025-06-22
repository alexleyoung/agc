package handlers

import "net/http"

func Setup(mux *http.ServeMux) {
	mux.HandleFunc("/auth/url", GetAuthURL)
	mux.HandleFunc("/auth/callback", OAuthCallback)
}
