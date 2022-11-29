package configuration

import "github.com/snasphysicist/ferp/v2/pkg/configuration/router"

// Configuration holds configuration for the entire application
type Configuration struct {
	Downstreams []Downstream `config:"downstreams"`
	HTTP        HTTP         `config:"http"`
	HTTPS       HTTPS        `config:"https"`
}

// Downstream represents a server that the proxy is providing access to
type Downstream struct {
	Target     string            `config:"target"`
	Protocol   string            `config:"protocol"`
	Host       string            `config:"host"`
	Port       uint16            `config:"port"`
	Base       string            `config:"base"`
	MapperData map[string]string `config:"path-mapper"`
	Mapper     pathMapper        `config:"-"`
}

// pathMapper is an object which can rewrite paths from the incoming request,
// to the downstream request, and can deserialise itself from a set of kv pairs
type pathMapper interface {
	Map(string) string
	From(map[string]string) error
}

// HTTP holds configuration for the HTTP proxy server
type HTTP struct {
	Port      uint16     `config:"port"`
	Redirects []Redirect `config:"redirect"`
	Incoming  []Incoming `config:"incoming"`
}

// HTTPS contains configuration for routes served by the proxy over HTTPS
type HTTPS struct {
	Port      uint16     `config:"port"`
	CertFile  string     `config:"cert-file"`
	KeyFile   string     `config:"key-file"`
	Redirects []Redirect `config:"redirects"`
	Incoming  []Incoming `config:"incoming"`
}

// Redirect configures the proxy to serve a redirect itself
type Redirect struct {
	From          string                `config:"from"`
	To            string                `config:"to"`
	Methods       []string              `config:"methods"`
	MethodRouters []router.MethodRouter `config:"-"` // populated after configuration load based on Methods
}

// Incoming represents a route that one of the proxy servers offers,
// and the target it proxies (the downstream)
type Incoming struct {
	Path          string                `config:"path"`
	Methods       []string              `config:"methods"`
	MethodRouters []router.MethodRouter `config:"-"` // populated after configuration load based on Methods
	Target        string                `config:"target"`
	Downstream    Downstream            `config:"-"` // populated after configuration load based on Target
}
