package cors

import (
	"net/http"
	"strings"
)

type Config struct {
	AllowedOrigins string
	AllowedMethods string
	AllowedHeaders string
	AllowedMaxAge  int
}

func New(cfg Config) func(http.Handler) http.Handler {
	allowedOrigins := strings.Split(cfg.AllowedOrigins, ",")
	allowedMethods := strings.Split(cfg.AllowedMethods, ",")
	allowedHeaders := strings.Split(cfg.AllowedHeaders, ",")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			originAllowed := false
			for _, allowed := range allowedOrigins {
				if strings.TrimSpace(allowed) == origin || strings.TrimSpace(allowed) == "*" {
					originAllowed = true
					break
				}
			}

			if originAllowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ", "))
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ", "))
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Max-Age", "3600")
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
