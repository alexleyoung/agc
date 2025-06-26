package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/alexleyoung/auto-gcal/internal/auth"
	"github.com/alexleyoung/auto-gcal/internal/db"
)

type contextKey string

const userKey contextKey = "user"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("agc_session")
		if err != nil {
			log.Print("Error getting session cookie:", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		session, err := db.GetSession(cookie.Value)
		if err != nil {
			log.Print("Error getting session:", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		expiresAt, err := time.Parse(time.RFC3339, session.ExpiresAt)
		if err != nil {
			log.Print("Error parsing session expiration time:", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if expiresAt.Before(time.Now()) {
			newToken, err := auth.RefreshToken(session)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			session.AccessToken = newToken.AccessToken
			session.ExpiresAt = newToken.Expiry.Format(time.RFC3339)
		}

		ctx := context.WithValue(r.Context(), userKey, session.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
