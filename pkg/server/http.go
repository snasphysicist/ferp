package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/snasphysicist/ferp/v2/pkg/configuration"
)

// HTTP sets up the HTTP proxy server, ready for starting
func HTTP(c configuration.HTTP) *http.Server {
	r := chi.NewRouter() // TODO: add logging middleware
	return &http.Server{Addr: fmt.Sprintf("%s:%d", host, c.Port), Handler: r}
}

// host is the host we serve on - always 0.0.0.0
const host = "0.0.0.0"
