package server

import (
	"fmt"
	"net/http"

	"github.com/snasphysicist/ferp/v2/pkg/configuration"
	"github.com/snasphysicist/ferp/v2/pkg/middleware"
	"github.com/snasphysicist/ferp/v2/pkg/server/forward"
	"github.com/snasphysicist/ferp/v2/pkg/server/redirect"
)

// HTTP sets up the HTTP proxy server, ready for starting
func HTTP(c configuration.HTTP) *http.Server {
	r := middleware.RouterWithDefaults()
	redirect.Configure(r, c.Redirects)
	forward.Configure(r, c.Incoming)
	return &http.Server{Addr: fmt.Sprintf("%s:%d", host, c.Port), Handler: r}
}

// host is the host we serve on - always 0.0.0.0
const host = "0.0.0.0"
