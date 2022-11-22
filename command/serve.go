package command

import (
	"os"
	"os/signal"

	"github.com/snasphysicist/ferp/v2/pkg/configuration"
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

// Serve TODO implement
func Serve(_ configuration.Configuration, stop chan struct{}) {
	shutdown := make(chan struct{})
	go shutDownOnSignalOrStop(shutdown, stop)
	<-shutdown
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
