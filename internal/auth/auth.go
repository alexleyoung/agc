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
	"time"

	"github.com/alexleyoung/auto-gcal/internal/db"
	"github.com/alexleyoung/auto-gcal/internal/types"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

var config *oauth2.Config

// Init initializes the OAuth2 client
func Init() {
	b, err := os.ReadFile("credentials.json")

	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	cfg, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope, calendar.CalendarEventsScope, "openid", "email", "profile")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	config = cfg
}

func GetClient(sessionID string) (*http.Client, error) {
	session, err := db.GetSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch user session:\n" + err.Error())
	}

	time, err := time.Parse(time.RFC3339, session.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse session expiration time:\n" + err.Error())
	}

	tok := &oauth2.Token{AccessToken: session.AccessToken, RefreshToken: session.RefreshToken, Expiry: time}

	return config.Client(context.Background(), tok), nil
}

func GetAuthURL() string {
	return config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
}

func Authenticate(r *http.Request) (types.User, types.Session, error) {
	// exchange code
	tok, err := config.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		return types.User{}, types.Session{}, err
	}

	// get id_token
	idTokenRaw, ok := tok.Extra("id_token").(string)
	if !ok {
		return types.User{}, types.Session{}, fmt.Errorf("Malformed token")
	}
	idToken, err := parseIDToken(idTokenRaw)
	if err != nil {
		return types.User{}, types.Session{}, fmt.Errorf("Failed to parse id_token: " + err.Error())
	}

	// create user and session as necessary
	user, err := db.CreateUser(idToken.Sub, idToken.Email, idToken.Name)
	if err != nil {
		return types.User{}, types.Session{}, fmt.Errorf("Failed to create user:\n" + err.Error())
	}
	session, err := db.CreateSession(idToken.Sub, tok.AccessToken, tok.RefreshToken, tok.Expiry)
	if err != nil {
		return types.User{}, types.Session{}, fmt.Errorf("Failed to create session:\n" + err.Error())
	}

	return user, session, nil
}

// takes JWT and parses payload for google sub and email
func parseIDToken(idToken string) (types.IDToken, error) {
	var claims types.IDToken

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
