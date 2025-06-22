package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/alexleyoung/auto-gcal/internal/auth"
	"github.com/alexleyoung/auto-gcal/internal/db"
	"golang.org/x/oauth2"
)

type UserInfo struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
}

// takes JWT and parses google sub from it
func extractUserInfoFromIDToken(idToken string) (UserInfo, error) {
	var claims UserInfo

	parts := strings.Split(idToken, ".")
	if len(parts) != 3 {
		return claims, fmt.Errorf("Malformed ID token")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return claims, err
	}

	if err = json.Unmarshal(payload, &claims); err != nil {
		return claims, err
	}

	return claims, nil
}

func getAuthURL(w http.ResponseWriter, r *http.Request) {
	config := auth.GetConfig()
	url := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusFound)
}

func oauthCallback(w http.ResponseWriter, r *http.Request) {
	config := auth.GetConfig()

	tok, err := config.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	idTokenRaw, ok := tok.Extra("id_token").(string)
	if !ok {
		http.Error(w, "No ID token found in token response", http.StatusInternalServerError)
	}

	info, err := extractUserInfoFromIDToken(idTokenRaw)
	if err != nil {
		http.Error(w, "Failed to extract user ID: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err = db.SaveToken(info.Sub, tok); err != nil {
		http.Error(w, "Failed to save token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Authenticated as user: " + info.Email))
}

func setupAuth(mux *http.ServeMux) {
	mux.HandleFunc("GET /auth", getAuthURL)
	mux.HandleFunc("GET /auth/callback", oauthCallback)
}
