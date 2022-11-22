package configuration

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// MethodRouter configures a route in the mux for a request with a specified method
// (or all methods) at the given path invoking the given handler
type MethodRouter interface {
	Route(r *chi.Mux, path string, handler http.HandlerFunc)
}
