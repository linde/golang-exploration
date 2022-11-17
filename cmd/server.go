package cmd

import (
	"log"
	"myapp/greeterserver"
	"myapp/grpcservice"

	"github.com/spf13/cobra"
)

func NewServerCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "example server for the greeter service",
		Run:   doServerRun,
	}
}

var serverCmd = NewServerCmd()

func init() {
	RootCmd.AddCommand(serverCmd)
}

func doServerRun(cmd *cobra.Command, args []string) {

	port, _ := cmd.Flags().GetInt("port")

	gs, err := grpcservice.NewServerFromPort(port)
	if err != nil {
		log.Fatalf("serverCmd.doServerRun failed to create server: %v", err)
	}
	helloServer := greeterserver.NewHelloServer()
	defer helloServer.Stop()
	serveErr := gs.Serve(helloServer)
	if serveErr != nil {
		log.Fatalf("serverCmd.doServerRun failed to serve: %v", err)
	}
}
