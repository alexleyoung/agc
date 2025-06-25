package auth

import (
	"context"
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

func GetClient(userID string) (*http.Client, error) {
	config := getSecretConfig()

	encTokString, err := db.GetUserToken(userID)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch user authentication token:\n" + err.Error())
	}

	tok, err := db.DecryptToken(encTokString)
	if err != nil {
		return nil, fmt.Errorf("Failed to decrypt token:\n", err.Error())
	}

	return config.Client(context.Background(), tok), nil
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

func VerifyAuthHeader(r *http.Request) (types.UserInfo, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return types.UserInfo{}, fmt.Errorf("Missing authorization header")
	}

	const bearerPrefix = "Bearer "
	if len(authHeader) <= len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		return types.UserInfo{}, fmt.Errorf("Invalid authorization header")
	}

	// takes JWT and parses payload for google sub and email
	parts := strings.Split(authHeader, ".")
	if len(parts) != 3 {
		return types.UserInfo{}, fmt.Errorf("Malformed authorization header")
	}
	payload, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return types.UserInfo{}, err
	}
	var claims types.UserInfo
	if err = json.Unmarshal(payload, &claims); err != nil {
		return types.UserInfo{}, err
	}

	return claims, nil
}
