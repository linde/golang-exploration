package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const DEFAULT_PORT = 10001

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   filepath.Base(os.Args[0]),
		Short: "implements a CLI for either a client or server of the greeter service",
	}

	return cmd
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

// TODO do ENV and file config support too

func init() {

}
