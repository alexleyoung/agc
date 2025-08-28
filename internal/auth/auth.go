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

	"github.com/alexleyoung/agc/internal/types"
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

	log.Print("Initialized OAuth2 client")
	config = cfg
}

func GetClient(session types.Session) (*http.Client, error) {
	time, err := time.Parse(time.RFC3339, session.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse session expiration time:\n%s", err.Error())
	}

	tok := &oauth2.Token{AccessToken: session.AccessToken, RefreshToken: session.RefreshToken, Expiry: time}

	return config.Client(context.Background(), tok), nil
}

func GetAuthURL() string {
	return config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
}

func Authenticate(r *http.Request, lookupTimezone func(session types.Session) (string, error)) (types.User, types.Session, error) {
	// exchange code
	tok, err := config.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		log.Print("Error exchanging code:", err)
		return types.User{}, types.Session{}, err
	}

	// get id_token
	idTokenRaw, ok := tok.Extra("id_token").(string)
	if !ok {
		log.Print("Malformed token")
		return types.User{}, types.Session{}, fmt.Errorf("Malformed token")
	}
	// id token
	_, err = parseIDToken(idTokenRaw)
	if err != nil {
		log.Print("Failed to parse id_token:", err)
		return types.User{}, types.Session{}, fmt.Errorf("Failed to parse id_token:\n%s", err.Error())
	}

	// create session
	// session, err := db.CreateSession(idToken.Sub, tok.AccessToken, tok.RefreshToken, tok.Expiry)
	// if err != nil {
	// 	log.Print("Failed to create session:", err)
	// 	return types.User{}, types.Session{}, fmt.Errorf("Failed to create session:\n%s", err.Error())
	// }
	//
	// user, err := db.GetUser(idToken.Sub)

	// if no user, create one
	// if user.UserID == "" {
	// 	timezone := ""
	// 	if lookupTimezone != nil {
	// 		timezone, err = lookupTimezone(session)
	// 		if err != nil {
	// 			log.Printf("Failed to look up timezone: %v", err)
	// 		}
	// 	}
	//
	// 	user, err = db.CreateUser(idToken.Sub, idToken.Email, idToken.Name, timezone)
	// 	if err != nil {
	// 		log.Print("Failed to create user:", err)
	// 		return types.User{}, types.Session{}, fmt.Errorf("Failed to create user:\n%s", err.Error())
	// 	}
	// }
	//
	// return user, session, nil
	return types.User{}, types.Session{}, nil
}

func RefreshToken(session types.Session) (*oauth2.Token, error) {
	expiry, err := time.Parse(time.RFC3339, session.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse session expiration time:\n%s", err.Error())
	}

	tok := &oauth2.Token{
		AccessToken:  session.AccessToken,
		RefreshToken: session.RefreshToken,
		Expiry:       expiry,
		TokenType:    "Bearer",
	}

	ts := config.TokenSource(context.Background(), tok)
	newToken, err := ts.Token()
	if err != nil {
		return nil, fmt.Errorf("token refresh failed: %w", err)
	}

	if newToken.AccessToken != session.AccessToken || newToken.Expiry != expiry {
		// err = db.UpdateSessionTokens(session.ID, newToken.AccessToken, newToken.Expiry)
		// if err != nil {
		// 	return nil, fmt.Errorf("failed to update session: %w", err)
		// }
	}

	return newToken, nil
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
