package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	return &cobra.Command{
		// TODO: get the name of the CLI command used from cobra
		Use:   "rpc-cmd-style",
		Short: "implements a CLI for either a client or server for the greeter service",
	}
}

// This represents the base command when called without any subcommands
var RootCmd = NewRootCmd()

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {

	// TODO do ENV and file config support too
	RootCmd.PersistentFlags().Int("port", 10001, "this is the service port")
}
