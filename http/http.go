package http

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"snas.pw/ferp/v2/configuration"
	"snas.pw/ferp/v2/routing"
)

const host = "0.0.0.0"

func Server(c configuration.HTTP) *http.Server {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	routing.SetUpRedirects(r, c.Redirects)
	routing.SetUpForwarding(r, c.Incoming)
	return &http.Server{Addr: fmt.Sprintf("%s:%d", host, c.Port), Handler: r}
}
