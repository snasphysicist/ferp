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
	log.L().Infof("%s request to %s", r.Method, r.URL.String())
	rr := wrap(w)
	next.ServeHTTP(&rr, r)
	log.L().Infof("Responded %d, wrote %d, took %s", rr.code, rr.written, time.Since(start))
}

// responseRecorder records some properties of the response which we want to log
type responseRecorder struct {
	code    int
	written int
	w       http.ResponseWriter
}

// wrap wraps the provided writer in a recorder
func wrap(w http.ResponseWriter) responseRecorder {
	return responseRecorder{w: w}
}

// WriteHeader captures the status code and forwards it to the wrapper writer
func (w *responseRecorder) WriteHeader(status int) {
	w.code = status
	w.w.WriteHeader(status)
}

// Write writes the bytes to the wrapper writer and records the number of bytes written
func (w *responseRecorder) Write(b []byte) (int, error) {
	n, err := w.w.Write(b)
	w.written = n
	return n, err
}

// Header just forwards to the wrapped writer's Header method
func (w *responseRecorder) Header() http.Header {
	return w.w.Header()
}
