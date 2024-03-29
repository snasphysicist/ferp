package integration

import (
	"context"
	"encoding/json"
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
	path   string
	method string
	rg     responseGenerator
}

// responseGenerator generates a response to a request in the mock
type responseGenerator func(r *http.Request) responseSpecification

type responseSpecification struct {
	status  int
	body    string
	headers http.Header
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
		hf := m.handler(r.rg)
		mux.Method(r.method, r.path, hf)
	}
	return mux
}

// handler creates a http.HandlerFunc which will always return a response
// with the given status code and whose body contains the given content
func (m mock) handler(rg responseGenerator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rs := rg(r)
		for k, vs := range rs.headers {
			for _, v := range vs {
				w.Header().Add(k, v)
			}
		}
		w.WriteHeader(rs.status)
		_, err := w.Write([]byte(rs.body))
		if err != nil {
			m.t.Errorf("Failed to write response content: %s", err)
		}
	}
}

// mockPorts are the ports on which mocks should be opened
func mockPorts() []uint16 {
	return []uint16{34543, 35753}
}

// setResponse simply returns the provided code and content for the response
func setResponse(code int, content string) responseGenerator {
	return func(*http.Request) responseSpecification {
		return responseSpecification{status: code, body: content, headers: make(http.Header)}
	}
}

// echoQueryParameters returns a 200 and writes
// the provided query parameters into the body as JSON
func echoQueryParameters() responseGenerator {
	return func(r *http.Request) responseSpecification {
		b, err := json.Marshal(r.URL.Query())
		if err != nil {
			panic(err)
		}
		return responseSpecification{status: 200, body: string(b), headers: make(http.Header)}
	}
}

// echoHeaders returns a 200 and writes
// the provided query parameters into the body as JSON
func echoHeaders() responseGenerator {
	return func(r *http.Request) responseSpecification {
		b, err := json.Marshal(r.Header)
		if err != nil {
			panic(err)
		}
		return responseSpecification{status: 200, body: string(b), headers: make(http.Header)}
	}
}
