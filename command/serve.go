package command

import (
	"context"
	"net/http"
	"os"
	"os/signal"

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
	go shutDownOnSignalOrStop(shutdown, stop)
	<-shutdown
}

// startInsecure starts the HTTP server according to its configuration, if needed
func startInsecure(c configuration.Configuration, shutdown <-chan struct{}) {
	if len(c.HTTP.Incoming) == 0 && len(c.HTTP.Redirects) == 0 {
		log.Infof("No HTTP routes or redirects configured, not starting HTTP")
		return
	}

	insecure := server.HTTP(c.HTTP)
	go func() {
		log.Infof("Starting http server on %s", insecure.Addr)
		err := insecure.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Errorf("http server stopped with %s", err)
			panic(err)
		}
	}()
	go shutDownGracefully(shutdown, insecure)
}

// shutDownOnSignal closes the shutdown channel if any OS signal or signal on stop is received
func shutDownOnSignalOrStop(shutdown chan<- struct{}, stop <-chan struct{}) {
	s := make(chan os.Signal, 1)
	signal.Notify(s)
	select {
	case <-s:
	case <-stop:
	}
	close(shutdown)
}

// shutDownGracefully attempts to shut down anything
// with a Shutdown method gracefully and logs any errors on shutdown
func shutDownGracefully(shutdown <-chan struct{}, s shutdowner) {
	<-shutdown
	if err := s.Shutdown(context.Background()); err != nil {
		log.Errorf("Failed to shut down %#v: %s", s, err)
	}
}

// shutdowner represents something which can be requested to shutdown gracefully
type shutdowner interface {
	Shutdown(context.Context) error
}
