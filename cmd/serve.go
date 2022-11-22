package cmd

import (
	"context"
	"io"
	"log"
	slhttp "net/http"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"snas.pw/ferp/v2/configuration"
	"snas.pw/ferp/v2/http"
	"snas.pw/ferp/v2/https"
)

// serveCommand sets up the command for starting the reverse proxy server
func serveCommand(c *configuration.Configuration) *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "start the reverse proxy & run until shutdown by a signal",
		Long:  "start the reverse proxy & run until shutdown by a signal",
		Run:   func(*cobra.Command, []string) { Serve(*c, make(chan struct{})) },
	}
}

// Serve starts the http and https reverse proxies according to the configuration
func Serve(c configuration.Configuration, stop chan struct{}) {
	shutdown := make(chan struct{})

	startInsecure(c, shutdown)
	startSecure(c, shutdown)

	go shutDownOnSignalOrStop(shutdown, stop)
	<-shutdown
}

// startInsecure starts the HTTP server according to its configuration, if needed
func startInsecure(c configuration.Configuration, shutdown <-chan struct{}) {
	if len(c.HTTP.Incoming) == 0 && len(c.HTTP.Redirects) == 0 {
		log.Printf("No HTTP routes or redirects configured, not starting HTTP")
		return
	}

	insecure := http.Server(c.HTTP)
	go func() {
		log.Printf("Starting http server on %s", insecure.Addr)
		err := insecure.ListenAndServe()
		if err != nil && err != slhttp.ErrServerClosed {
			log.Printf("http server stopped with %s", err)
			panic(err)
		}
	}()
	go shutDownGracefully(shutdown, insecure)
}

// startSecure starts the HTTPS server according to its configuration, if needed
func startSecure(c configuration.Configuration, shutdown <-chan struct{}) {
	if len(c.HTTPS.Incoming) == 0 && len(c.HTTPS.Redirects) == 0 {
		log.Printf("No HTTPS routes or redirects configured, not starting HTTPS")
		return
	}

	ensureExists(c.HTTPS.CertFile)
	ensureExists(c.HTTPS.KeyFile)
	secure := https.Server(c.HTTPS)
	go func() {
		log.Printf("Starting https server on %s", secure.Addr)
		err := secure.ListenAndServeTLS(c.HTTPS.CertFile, c.HTTPS.KeyFile)
		if err != nil && err != slhttp.ErrServerClosed {
			log.Printf("https server stopped with %s", err)
			panic(err)
		}
	}()
	go shutDownGracefully(shutdown, secure)
}

// ensureExists panics if a file at the given path does not exist
func ensureExists(path string) {
	f, err := os.Open(path)
	defer closeLoggingErrors(f)
	if err != nil {
		log.Panicf("Failed to open %s", path)
	}
}

// shutDownGracefully attempts to shut down anything
// with a Shutdown method gracefully and logs any errors on shutdown
func shutDownGracefully(shutdown <-chan struct{}, s shutdowner) {
	<-shutdown
	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("Failed to shut down %#v: %s", s, err)
	}
}

// shutdowner represents something which can be requested to shutdown gracefully
type shutdowner interface {
	Shutdown(context.Context) error
}

// closeLoggingErrors closes a closeable and logs any errors encountered on close
func closeLoggingErrors(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Printf("Failed to close %#v: %s", c, err)
	}
}

// shutDownOnSignal closes the provided channel if any signal is received
func shutDownOnSignalOrStop(shutdown chan<- struct{}, stop <-chan struct{}) {
	s := make(chan os.Signal, 1)
	signal.Notify(s)
	select {
	case <-s:
	case <-stop:
	}
	close(shutdown)
}
