package router

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// RouterFor finds an appropriate router to set up routes for the given method
func RouterFor(method string) (MethodRouter, error) {
	mr, ok := MethodRouters()[method]
	if !ok {
		return nil, fmt.Errorf("method '%s' is not supported", method)
	}
	return mr, nil
}

// MethodRouters lists all method types that are recognised by this application
// and the method routers for them
func MethodRouters() map[string]MethodRouter {
	return map[string]MethodRouter{
		http.MethodConnect: connect{},
		http.MethodDelete:  delete{},
		http.MethodGet:     get{},
		http.MethodHead:    head{},
		http.MethodOptions: options{},
		http.MethodPatch:   patch{},
		http.MethodPost:    post{},
		http.MethodPut:     put{},
		http.MethodTrace:   trace{},
		"*":                all{},
	}
}

// MethodRouter configures a route in the mux for a request with a specified method
// (or, in one special case, all methods) at the given path invoking the given handler
type MethodRouter interface {
	Route(r *chi.Mux, path string, handler http.HandlerFunc)
}

// Each MethodRouter implementation is named after the method for which it sets up routes
// and just passes through to a call of the same named method on the router,
// except for all which calls Handle (which routes any method).

type connect struct{}

func (connect) Route(r *chi.Mux, path string, handler http.HandlerFunc) {
	r.Connect(path, handler)
}

type delete struct{}

func (delete) Route(r *chi.Mux, path string, handler http.HandlerFunc) {
	r.Delete(path, handler)
}

type get struct{}

func (get) Route(r *chi.Mux, path string, handler http.HandlerFunc) {
	r.Get(path, handler)
}

type head struct{}

func (head) Route(r *chi.Mux, path string, handler http.HandlerFunc) {
	r.Head(path, handler)
}

type options struct{}

func (options) Route(r *chi.Mux, path string, handler http.HandlerFunc) {
	r.Options(path, handler)
}

type patch struct{}

func (patch) Route(r *chi.Mux, path string, handler http.HandlerFunc) {
	r.Patch(path, handler)
}

type post struct{}

func (post) Route(r *chi.Mux, path string, handler http.HandlerFunc) {
	r.Post(path, handler)
}

type put struct{}

func (put) Route(r *chi.Mux, path string, handler http.HandlerFunc) {
	r.Put(path, handler)
}

type trace struct{}

func (trace) Route(r *chi.Mux, path string, handler http.HandlerFunc) {
	r.Trace(path, handler)
}

type all struct{}

func (all) Route(r *chi.Mux, path string, handler http.HandlerFunc) {
	r.Handle(path, handler)
}
