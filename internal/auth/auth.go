package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/alexleyoung/auto-gcal/internal/db"
	"github.com/alexleyoung/auto-gcal/internal/types"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

type UserInfo struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
}

func GetClient() *http.Client {
	// config := getSecretConfig()
	return nil
}

func getSecretConfig() *oauth2.Config {
	b, err := os.ReadFile("credentials.json")

	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope, calendar.CalendarEventsScope, "openid", "email")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	return config
}

func GetAuthURL() string {
	config := getSecretConfig()
	return config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
}

func Authenticate(r *http.Request) (types.UserInfo, error) {
	config := getSecretConfig()

	tok, err := config.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		return types.UserInfo{}, err
	}

	idTokenRaw, ok := tok.Extra("id_token").(string)
	if !ok {
		return types.UserInfo{}, fmt.Errorf("Malformed token")
	}

	info, err := extractUserInfoFromIDToken(idTokenRaw)
	if err != nil {
		return types.UserInfo{}, fmt.Errorf("Failed to parse id_token: " + err.Error())
	}

	if err = db.SaveToken(info, tok); err != nil {
		return types.UserInfo{}, fmt.Errorf("Failed to save token: " + err.Error())
	}

	return info, nil
}

// takes JWT and parses payload for google sub and email
func extractUserInfoFromIDToken(idToken string) (types.UserInfo, error) {
	var claims types.UserInfo

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
