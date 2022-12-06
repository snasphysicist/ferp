package server

import (
	"fmt"
	"net/http"

	"github.com/snasphysicist/ferp/v2/pkg/configuration"
	"github.com/snasphysicist/ferp/v2/pkg/middleware"
	"github.com/snasphysicist/ferp/v2/pkg/server/forward"
	"github.com/snasphysicist/ferp/v2/pkg/server/redirect"
)

// HTTPS sets up the HTTPS proxy server, returning it ready to serve
func HTTPS(c configuration.HTTPS) *http.Server {
	r := middleware.RouterWithDefaults()
	redirect.Configure(r, c.Redirects)
	forward.Configure(r, c.Incoming)
	return &http.Server{Addr: fmt.Sprintf("0.0.0.0:%d", c.Port), Handler: r}
}
