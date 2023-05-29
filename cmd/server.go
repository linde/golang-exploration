package cmd

import (
	"log"
	"myapp/grpcservice"
	"myapp/helloserver"
	"myapp/restserver"

	"github.com/spf13/cobra"
)

func NewServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "example server for the greeter service",
		Run:   doServerRun,
	}

	// TODO document and e2e test the rest port
	cmd.Flags().IntVarP(&restPort, "rest", "r", -1, "port to use to also enable the rest gateway")
	return cmd

}

var serverCmd = NewServerCmd()
var restPort int

func init() {
	RootCmd.AddCommand(serverCmd)
}

func doServerRun(cmd *cobra.Command, args []string) {

	rpcPort, _ := cmd.Flags().GetInt("port") // this references the root param --port

	gs, err := grpcservice.NewServerFromPort(rpcPort)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}
	helloServer := helloserver.NewHelloServer()
	defer helloServer.Stop()

	// if restPort is configured, add a rest gateway using the port
	if restPort >= 0 {
		rpcAddr, err := gs.GetServiceTCPAddr()
		if err != nil {
			log.Fatalf("error getting RPC service address: %s", err)
		}
		go restserver.NewRestGateway(restPort, rpcAddr)
	}

	serveErr := gs.Serve(helloServer)
	if serveErr != nil {
		log.Fatalf("serverCmd.doServerRun failed to serve: %v", err)
	}

}
