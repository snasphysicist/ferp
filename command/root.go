package command

import (
	"github.com/snasphysicist/ferp/v2/pkg/log"
	"github.com/spf13/cobra"
)

// Execute is the main entry point for the whole application
func Execute() {
	defer initialiseLogging()()
	root := cobra.Command{
		Use:   "ferp",
		Short: "fabulously easy reverse proxy",
		Long:  "a super easy to use reverse proxy, that supports http & https incoming, and multiple downstream services",
	}
	var path string
	root.PersistentFlags().StringVar(&path, "configuration-file", "", "path to the proxy configuration file")
	_ = root.Execute()
}

// initialiseLogging initialises the logging package
func initialiseLogging() func() {
	flush, err := log.Initialise()
	if err != nil {
		panic(err)
	}
	return flush
}
