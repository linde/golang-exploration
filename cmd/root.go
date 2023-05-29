package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   filepath.Base(os.Args[0]),
		Short: "implements a CLI for either a client or server of the greeter service",
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

// TODO do ENV and file config support too

func init() {

	RootCmd.PersistentFlags().IntP("port", "p", 10001, "grpc server port to use or connect to")
}
