package command

import (
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
func Serve(_ configuration.Configuration, _ chan struct{}) {}
