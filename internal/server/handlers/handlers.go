package handlers

import "net/http"

func Init(mux *http.ServeMux) {
	setupAI(mux)
	setupAuth(mux)
}
