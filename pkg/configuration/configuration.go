package configuration

import "github.com/snasphysicist/ferp/v2/pkg/configuration/router"

// Configuration holds configuration for the entire application
type Configuration struct {
	HTTP HTTP `config:"http"`
}

// HTTP holds configuration for the HTTP proxy server
type HTTP struct {
	Port      uint16     `config:"port"`
	Redirects []Redirect `config:"redirect"`
}

// Redirect configures the proxy to serve a redirect itself
type Redirect struct {
	From          string                `config:"from"`
	To            string                `config:"to"`
	Methods       []string              `config:"methods"`
	MethodRouters []router.MethodRouter `config:"-"` // populated after configuration load based on Methods
}
