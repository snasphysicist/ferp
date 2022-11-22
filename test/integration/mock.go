package integration

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
)

type mock struct {
	t      *testing.T
	port   uint16
	routes []route
}

type route struct {
	path    string
	method  string
	code    int
	content string
}

func (m mock) start() func() {
	s := http.Server{
		Addr:    fmt.Sprintf("localhost:%d", m.port),
		Handler: m.configureRoutes(),
	}
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			m.t.Errorf("Failed to start mock server: %s", err)
		}
	}()
	return func() {
		if err := s.Shutdown(context.Background()); err != nil {
			m.t.Errorf("Failed to shut down mock server: %s", err)
		}
	}
}

func (m mock) configureRoutes() *chi.Mux {
	mux := chi.NewMux()
	for _, r := range m.routes {
		hf := m.handler(r.code, r.content)
		mux.Method(r.method, r.path, hf)
	}
	return mux
}

func (m mock) handler(code int, content string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		_, err := w.Write([]byte(content))
		if err != nil {
			m.t.Errorf("Failed to write response content: %s", err)
		}
	}
}
