package cmd

import (
	"log"

	greeterserver "myapp/greeterserver"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "example server for the greeter service",
	Run:   doServerRun,
}

func init() {
	RootCmd.AddCommand(serverCmd)
}

func doServerRun(cmd *cobra.Command, args []string) {

	port, _ := cmd.Flags().GetInt("port")

	err := greeterserver.ServePort(port)
	if err != nil {
		log.Fatalf("serverCmd failed to serve: %v", err)
	}
}
