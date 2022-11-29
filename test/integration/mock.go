package integration

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
)

// mock is a mock HTTP server used for integration testing.
// It serves preconfigured responses on a set of preconfigured routes.
// It will fail the test if any part of the setup/serving fails unexpectedly.
type mock struct {
	t      *testing.T
	port   uint16
	routes []route
}

// route is a route offered by the mock, with the code/content that will be returned from that route
type route struct {
	path    string
	method  string
	code    int
	content string
}

// start starts the mock server in a new goroutine, and returns a function
// which should be deferred from the test to shutdown that server.
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

// configureRoutes sets up a router according to the mock's configured routes
func (m mock) configureRoutes() *chi.Mux {
	mux := chi.NewMux()
	for _, r := range m.routes {
		hf := m.handler(r.code, r.content)
		mux.Method(r.method, r.path, hf)
	}
	return mux
}

// handler creates a http.HandlerFunc which will always return a response
// with the given status code and whose body contains the given content
func (m mock) handler(code int, content string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		_, err := w.Write([]byte(content))
		if err != nil {
			m.t.Errorf("Failed to write response content: %s", err)
		}
	}
}

// mockPorts are the ports on which mocks should be opened
func mockPorts() []uint16 {
	return []uint16{34543, 35753}
}
