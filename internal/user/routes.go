package user

import (
	"encoding/json"
	"net/http"
)

func Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /users/me", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode("hmm")
	})
}
