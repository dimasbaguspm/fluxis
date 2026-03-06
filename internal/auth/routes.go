package auth

import (
	"encoding/json"
	"net/http"
)

func Routes(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/register", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode("auth - register")
	})
	mux.HandleFunc("POST /auth/login", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode("auth - login")
	})
	mux.HandleFunc("POST /auth/refresh", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode("auth - refresh")
	})
}
