package command

import "github.com/spf13/cobra"

// Execute is the main entry point for the whole application
func Execute() {
	root := cobra.Command{
		Use:   "ferp",
		Short: "fabulously easy reverse proxy",
		Long:  "a super easy to use reverse proxy, that supports http & https incoming, and multiple downstream services",
	}
	var path string
	root.PersistentFlags().StringVar(&path, "configuration-file", "", "path to the proxy configuration file")
	_ = root.Execute()
}
