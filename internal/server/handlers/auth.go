package handlers

import (
	"net/http"

	"github.com/alexleyoung/agc/internal/auth"
	"github.com/alexleyoung/agc/internal/calendar"
	"github.com/alexleyoung/agc/internal/types"
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
	user, session, err := auth.Authenticate(r, func(session types.Session) (string, error) {
		cal, err := calendar.GetCalendar(r.Context(), session, "primary")
		if err != nil {
			return "", err
		}
		return cal.TimeZone, nil
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "agc_session",
		Value:    session.ID,
		HttpOnly: true,
		Secure:   false, // TEMP: CHANGE TO TRUE FOR PROD!!!!
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	w.Write([]byte("Authenticated as user: " + user.Email))
}
