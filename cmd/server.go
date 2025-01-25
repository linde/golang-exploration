package cmd

import (
	"fmt"
	"myapp/grpcservice"
	"myapp/helloserver"
	"myapp/restserver"

	"github.com/spf13/cobra"
)

var serverCmd = NewServerCmd()
var restPort, rpcPort int

func NewServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "example server for the greeter service",
		RunE:  doServerRun,
	}

	// TODO document and e2e test the rest port

	cmd.Flags().IntVarP(&rpcPort, "port", "p", DEFAULT_PORT, "rpcserver port")
	cmd.Flags().IntVarP(&restPort, "rest", "r", -1, "port to use to also enable the rest gateway")

	return cmd
}

func init() {
	RootCmd.AddCommand(serverCmd)
}

func doServerRun(cmd *cobra.Command, args []string) error {

	gs, err := grpcservice.NewServerFromPort(rpcPort)
	if err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "failed to create server: %v", err)
		return err
	}
	helloServer := helloserver.NewHelloServer()
	defer helloServer.Stop()

	// if restPort is configured, add a rest gateway using the port
	if restPort >= 0 {
		rpcAddr, err := gs.GetServiceTCPAddr()
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "error getting RPC service address: %s", err)
			return err
		}
		rgw := restserver.NewRestGateway(restPort, rpcAddr)
		go rgw.Serve()

	}

	serveErr := gs.Serve(helloServer)
	if serveErr != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "error in grpc server Serve(): %v", err)
		return err
	}

	return nil
}
