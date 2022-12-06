package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/snasphysicist/ferp/v2/pkg/log"
)

// RouterWithDefaults returns a new router with the default
// middlewares for the proxy servers attached
func RouterWithDefaults() *chi.Mux {
	r := chi.NewRouter()
	r.Use(logging)
	return r
}

// logging is a logging middleware which uses the logger from this service
func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logRequest(w, r, next)
	})
}

// logRequest logs some information about the request (method, url, timing, ...)
func logRequest(w http.ResponseWriter, r *http.Request, next http.Handler) {
	start := time.Now()
	log.Infof("%s request to %s", r.Method, r.URL.String())
	next.ServeHTTP(w, r)
	log.Infof("Took %s", time.Since(start))
}
