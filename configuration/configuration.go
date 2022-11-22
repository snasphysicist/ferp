package configuration

import (
	"fmt"
	"log"
	"net/http"
	gopath "path"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"snas.pw/ferp/v2/functional"
	"snas.pw/ferp/v2/mapper"
)

// Configuration contains all configuration for the application
type Configuration struct {
	Downstreams []Downstream `config:"downstream"`
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

// pathMapper is an object which can rewrite paths
// from the incoming request to the downstream request
type pathMapper interface {
	Map(string) string
	From(map[string]string) error
}

// HTTP contains configuration for routes served
// by the proxy over HTTP (i.e. not HTTPS)
type HTTP struct {
	Port      uint16     `config:"port"`
	Redirects []Redirect `config:"redirects"`
	Incoming  []Incoming `config:"incoming"`
}

// Incoming represents a route that one of the proxy servers
// offers, and the target it proxies (the downstream)
type Incoming struct {
	Path          string         `config:"path"`
	Methods       []string       `config:"methods"`
	MethodRouters []MethodRouter `config:"-"`
	Target        string         `config:"target"`
	Downstream    Downstream     `config:"downstream"`
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
	From          string         `config:"from"`
	To            string         `config:"to"`
	Methods       []string       `config:"methods"`
	MethodRouters []MethodRouter `config:"-"`
}

// Load attempts to load and validate the server configuration from the given path
func Load(path string) (Configuration, error) {
	err := readInConfiguration(path)
	if err != nil {
		return Configuration{}, err
	}
	c := Configuration{}
	err = viper.Unmarshal(&c, func(c *mapstructure.DecoderConfig) { c.TagName = "config" })
	if err != nil {
		log.Printf("Failed to deserialise configuration: %s", err)
		return Configuration{}, err
	}
	errs := functional.Filter([]error{loadPathMappers(&c), validateDownstreams(&c), validateAllMethods(c)}, notNil)
	errStr := maps(errs, func(e error) string { return e.Error() })
	if len(errStr) > 0 {
		allErrs := strings.Join(errStr, "; ")
		log.Printf("Failed to validate configuration: %s", allErrs)
		return Configuration{}, fmt.Errorf("invalid configuration: %s", allErrs)
	}
	assignRoutersForMethods(&c)
	log.Printf("Loaded configuration: %#v", c)
	return c, nil
}

// readInConfiguration configures viper to target the configuration file and attempts to read it in
func readInConfiguration(path string) error {
	directory := filepath.Dir(path)
	viper.AddConfigPath(directory)
	name := strings.TrimSuffix(gopath.Base(path), gopath.Ext(path))
	viper.SetConfigName(name)
	log.Printf("Attempting to load configuration from directory %s, name %s",
		directory, name)
	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Failed to read in configuration: %s", err)
	}
	return err
}

// validateDownstreams ensures that all downstreams specified in the server configurations exist
func validateDownstreams(c *Configuration) error {
	errs := make([]error, 0)
	for i := range c.HTTP.Incoming {
		err := findDownstream(c.Downstreams, &c.HTTP.Incoming[i])
		errs = append(errs, err)
	}
	for i := range c.HTTPS.Incoming {
		err := findDownstream(c.Downstreams, &c.HTTPS.Incoming[i])
		errs = append(errs, err)
	}
	errs = functional.Filter(errs, notNil)
	errStr := maps(errs, func(e error) string { return e.Error() })
	if len(errs) != 0 {
		return fmt.Errorf("invalid downstream settings: %+v", strings.Join(errStr, ", "))
	}
	return nil
}

// validateAllMethods ensures that all configured incoming HTTP methods are valid options
func validateAllMethods(c Configuration) error {
	httpErrs := functional.Filter(
		maps(c.HTTP.Incoming, func(i Incoming) error { return validateMethods(c, i.Methods, i) }),
		notNil)
	httpsErrs := functional.Filter(
		maps(c.HTTPS.Incoming, func(i Incoming) error { return validateMethods(c, i.Methods, i) }),
		notNil)
	httpRedirectErrs := functional.Filter(
		maps(c.HTTP.Redirects, func(r Redirect) error { return validateMethods(c, r.Methods, r) }),
		notNil)
	httpsRedirectErrs := functional.Filter(
		maps(c.HTTPS.Redirects, func(r Redirect) error { return validateMethods(c, r.Methods, r) }),
		notNil)
	errs := concat(httpErrs, httpsErrs, httpRedirectErrs, httpsRedirectErrs)
	errStr := maps(errs, func(e error) string { return e.Error() })
	if len(errs) != 0 {
		return fmt.Errorf("invalid methods settings: %+v", strings.Join(errStr, ", "))
	}
	return nil
}

// findDownstream attempts to find the downstream in the configuration for the given incoming route
func findDownstream(d []Downstream, i *Incoming) error {
	matches := functional.Filter(d, func(d Downstream) bool { return d.Target == i.Target })
	if len(matches) != 1 {
		return fmt.Errorf(
			"invalid downstream target %s from %#v in configuration %#v",
			i.Target, i, d)
	}
	i.Downstream = matches[0]
	return nil
}

// validateMethods validates that the configured HTTP methods in the given incoming are valid options
func validateMethods(c Configuration, ms []string, source interface{}) error {
	errs := functional.Filter(maps(ms, func(m string) error { return validateMethod(c, m) }), notNil)
	errStr := maps(errs, func(e error) string { return e.Error() })
	if len(errs) != 0 {
		return fmt.Errorf("invalid methods in %#v: %s", source, strings.Join(errStr, ", "))
	}
	return nil
}

// validateMethod validates that the provided HTTP methods is a valid option
func validateMethod(c Configuration, m string) error {
	if !contains(validMethods(), m) {
		return fmt.Errorf("invalid method %s", m)
	}
	return nil
}

// maps implements slice mapping functionality
func maps[T any, U any](s []T, f func(T) U) []U {
	m := make([]U, len(s))
	for i, e := range s {
		m[i] = f(e)
	}
	return m
}

// contains returns true iff something equal to e can be found as an element of s
func contains[T comparable](s []T, e T) bool {
	return len(functional.Filter(s, func(m T) bool { return m == e })) > 0
}

// notNil returns true iff e is not nil
func notNil(e error) bool {
	return e != nil
}

// concat returns a slice containing all elements of all provided slices,
// preserving element order both between and within the input slices
func concat[T any](ss ...[]T) []T {
	all := make([]T, 0)
	for _, s := range ss {
		all = append(all, s...)
	}
	return all
}

// validMethods lists all valid HTTP method options in the incoming configurations.
// The special value "ALL" means proxy requests for all method types.
func validMethods() []string {
	ms := make([]string, 0)
	for m := range methodRouters() {
		ms = append(ms, m)
	}
	return ms
}

// loadPathMappers attempts to find and instantiate a path mapper
// for all downstream path mapper configurations
// TODO: check and return all errors, not just first one
func loadPathMappers(c *Configuration) error {
	for i := range c.Downstreams {
		m, err := loadPathMapper(c.Downstreams[i].MapperData)
		if err != nil {
			return err
		}
		c.Downstreams[i].Mapper = m
	}
	return nil
}

// loadPathMapper attempts to find and instantiate a path mapper
// matching the provided configuration map
// TODO: collect and return all errors
// TODO: all extensibility with custom path mappers added at "runtime"
func loadPathMapper(c map[string]string) (pathMapper, error) {
	mappers := []pathMapper{&mapper.Passthrough{}, &mapper.RemovePrefix{}}
	for _, m := range mappers {
		err := m.From(c)
		if err == nil {
			return m, nil
		}
	}
	return nil, fmt.Errorf("no mapper matching configuration %#v", c)
}

// MethodRouter configures a route in the mux for a request with a specified method
// (or all methods) at the given path invoking the given handler
type MethodRouter func(r *chi.Mux, path string, handler http.HandlerFunc)

// assignRoutersForMethods sets up a MethodRouter for each method in each incoming
func assignRoutersForMethods(c *Configuration) {
	for i := range c.HTTP.Incoming {
		for _, m := range c.HTTP.Incoming[i].Methods {
			c.HTTP.Incoming[i].MethodRouters = append(
				c.HTTP.Incoming[i].MethodRouters, routerFor(m))
			log.Printf("For %#v: %#v", c.HTTP.Incoming[i], c.HTTP.Incoming[i].MethodRouters)
		}
	}
	for i := range c.HTTPS.Incoming {
		for _, m := range c.HTTPS.Incoming[i].Methods {
			c.HTTPS.Incoming[i].MethodRouters = append(
				c.HTTPS.Incoming[i].MethodRouters, routerFor(m))
		}
	}
	for i := range c.HTTP.Redirects {
		for _, m := range c.HTTP.Redirects[i].Methods {
			c.HTTP.Redirects[i].MethodRouters = append(
				c.HTTP.Redirects[i].MethodRouters, routerFor(m))
		}
	}
	for i := range c.HTTPS.Redirects {
		for _, m := range c.HTTPS.Redirects[i].Methods {
			c.HTTPS.Redirects[i].MethodRouters = append(
				c.HTTPS.Redirects[i].MethodRouters, routerFor(m))
		}
	}
}

// routerFor finds an appropriate router to set up routes for the given method
func routerFor(method string) MethodRouter {
	mr, ok := methodRouters()[method]
	if !ok {
		log.Panicf("unknown method %s", method)
	}
	return mr
}

// methodRouters lists all routers for each method that's valid in configuration
func methodRouters() map[string]MethodRouter {
	return map[string]MethodRouter{
		http.MethodConnect: func(r *chi.Mux, path string, handler http.HandlerFunc) { r.Connect(path, handler) },
		http.MethodDelete:  func(r *chi.Mux, path string, handler http.HandlerFunc) { r.Delete(path, handler) },
		http.MethodGet:     func(r *chi.Mux, path string, handler http.HandlerFunc) { r.Get(path, handler) },
		http.MethodHead:    func(r *chi.Mux, path string, handler http.HandlerFunc) { r.Head(path, handler) },
		http.MethodOptions: func(r *chi.Mux, path string, handler http.HandlerFunc) { r.Options(path, handler) },
		http.MethodPatch:   func(r *chi.Mux, path string, handler http.HandlerFunc) { r.Patch(path, handler) },
		http.MethodPost:    func(r *chi.Mux, path string, handler http.HandlerFunc) { r.Post(path, handler) },
		http.MethodPut:     func(r *chi.Mux, path string, handler http.HandlerFunc) { r.Put(path, handler) },
		http.MethodTrace:   func(r *chi.Mux, path string, handler http.HandlerFunc) { r.Trace(path, handler) },
		"ALL":              func(r *chi.Mux, path string, handler http.HandlerFunc) { r.Handle(path, handler) },
	}
}
