package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"snas.pw/ferp/v2/configuration"
)

// Execute is the main entry point for the application
func Execute() {
	root := cobra.Command{
		Use:   "ferp",
		Short: "fabulously easy reverse proxy",
		Long:  "a super easy to use reverse proxy, that supports http & https incoming, and multiple downstream services",
	}
	var path string
	root.PersistentFlags().StringVar(&path, "configuration-file", "", "path to the proxy configuration file")
	var c configuration.Configuration
	cobra.OnInitialize(func() { loadConfiguration(&path, &c) })
	root.AddCommand(serveCommand(&c))
	_ = root.Execute()
}

// loadConfiguration loads the configuration from the passed path into the passed struct
func loadConfiguration(path *string, c *configuration.Configuration) {
	cl, err := configuration.Load(*path)
	if err != nil {
		log.Panicf("Failed to load configuration: %s", err)
	}
	*c = cl
}
