package redirect

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/snasphysicist/ferp/v2/pkg/configuration"
	"github.com/snasphysicist/ferp/v2/pkg/log"
)

// Configure sets up all redirects specified in configuration on the provided router
func Configure(r *chi.Mux, rds []configuration.Redirect) {
	for _, rd := range rds {
		for _, mr := range rd.MethodRouters {
			mr.Route(r, rd.From, redirector(rd.To))
			log.L().Infof("Configuring redirect from '%s' to '%s' with %#v", rd.From, rd.To, mr)
		}
	}
}

// redirector creates a HTTP handler which returns a redirect (302) to the provided URL
func redirector(to string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("location", to)
		w.WriteHeader(http.StatusFound)
	}
}
