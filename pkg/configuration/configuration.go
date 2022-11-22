package configuration

// Configuration holds configuration for the entire application
type Configuration struct {
	HTTP HTTP `config:"http"`
}

// HTTP holds configuration for the HTTP proxy server
type HTTP struct {
	Port uint16 `config:"port"`
}
