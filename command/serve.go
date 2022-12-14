package command

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/snasphysicist/ferp/v2/pkg/configuration"
	"github.com/snasphysicist/ferp/v2/pkg/log"
	"github.com/snasphysicist/ferp/v2/pkg/server"
	"github.com/spf13/cobra"
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

// Serve starts the HTTP and HTTPS proxy servers, and runs them until signalled
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
		log.L().Infof("No HTTP routes or redirects configured, not starting HTTP")
		return
	}

	insecure := server.HTTP(c.HTTP)
	go func() {
		log.L().Infof("Starting http server on %s", insecure.Addr)
		err := insecure.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.L().Errorf("http server stopped with %s", err)
			panic(err)
		}
	}()
	go shutDownGracefully(shutdown, insecure)
}

// startSecure starts the HTTPS server according to its configuration, if needed
func startSecure(c configuration.Configuration, shutdown <-chan struct{}) {
	if len(c.HTTPS.Incoming) == 0 && len(c.HTTPS.Redirects) == 0 {
		log.L().Infof("No HTTPS routes or redirects configured, not starting HTTPS")
		return
	}

	ensureExists(c.HTTPS.CertFile)
	ensureExists(c.HTTPS.KeyFile)
	secure := server.HTTPS(c.HTTPS)
	go func() {
		log.L().Infof("Starting https server on %s", secure.Addr)
		err := secure.ListenAndServeTLS(c.HTTPS.CertFile, c.HTTPS.KeyFile)
		if err != nil && err != http.ErrServerClosed {
			log.L().Errorf("https server stopped with %s", err)
			panic(err)
		}
	}()
	go shutDownGracefully(shutdown, secure)
}

// shutDownOnSignal closes the shutdown channel if any OS signal or signal on stop is received
func shutDownOnSignalOrStop(shutdown chan<- struct{}, stop <-chan struct{}) {
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt, syscall.SIGHUP, syscall.SIGINT)
	select {
	case sgn := <-s:
		log.L().Infof("Received system shutdown signal %v", sgn)
	case <-stop:
		log.L().Infof("Received internal shutdown signal")
	}
	close(shutdown)
}

// shutDownGracefully attempts to shut down anything
// with a Shutdown method gracefully and logs any errors on shutdown
func shutDownGracefully(shutdown <-chan struct{}, s shutdowner) {
	<-shutdown
	if err := s.Shutdown(context.Background()); err != nil {
		log.L().Errorf("Failed to shut down %#v: %s", s, err)
	}
}

// shutdowner represents something which can be requested to shutdown gracefully
type shutdowner interface {
	Shutdown(context.Context) error
}

// ensureExists panics if a file at the given path does not exist
func ensureExists(path string) {
	f, err := os.Open(path)
	defer closeLoggingErrors(f)
	if err != nil {
		log.L().Errorf("Failed to open %s", path)
		panic(err)
	}
}

// closeLoggingErrors closes a closeable and logs any errors encountered on close
func closeLoggingErrors(c io.Closer) {
	if err := c.Close(); err != nil {
		log.L().Errorf("Failed to close %#v: %s", c, err)
	}
}
