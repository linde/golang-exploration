package cmd

import (
	"errors"
	"fmt"
	"myapp/grpcservice"
	"myapp/helloserver"
	"myapp/restserver"
	"net"
	"time"

	"github.com/spf13/cobra"
)

var serverCmd = NewServerCmd()
var serverCmdInfo *ServerCmdInfo
var requestedRestPort, requestedRpcPort int

// TODO should really make a new struct to bring the info and cmd together
// TODO should probably reason with Addr's not ports maybe
type ServerCmdInfo struct {
	rpcServingPort int
	rpcStopFunc    func()

	restInitialized bool
	restServingAddr net.Addr
	restStopFunc    func()
}

func ServerCommandClose() {
	if serverCmdInfo == nil {
		return
	}

	serverCmdInfo.rpcStopFunc()
	if serverCmdInfo.restInitialized {
		serverCmdInfo.restStopFunc()
	}
	serverCmdInfo = nil
}

func NewServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "example server for the greeter service",
		RunE:  doServerRun,
	}

	// TODO document and e2e test the rest port

	cmd.Flags().IntVarP(&requestedRpcPort, "port", "p", DEFAULT_PORT, "rpcserver port")
	cmd.Flags().IntVarP(&requestedRestPort, "rest", "r", -1, "port to use to also enable the rest gateway")

	return cmd
}

func init() {
	RootCmd.AddCommand(serverCmd)
}

func doServerRun(cmd *cobra.Command, args []string) error {

	gs, err := grpcservice.NewServerFromPort(requestedRpcPort)
	if err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "failed to create server: %v", err)
		return err
	}
	helloServer := helloserver.NewHelloServer()
	defer helloServer.Stop()

	if rpcServingPort, err := gs.GetServicePort(); err == nil {
		serverCmdInfo = &ServerCmdInfo{
			rpcServingPort:  rpcServingPort,
			rpcStopFunc:     helloServer.Stop,
			restInitialized: false,
		}
	}

	// if restPort is configured, add a rest gateway using the port
	if requestedRestPort >= 0 {
		rpcAddr, err := gs.GetServiceTCPAddr()
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "error getting RPC service address: %s", err)
			return err
		}
		rgw := restserver.NewRestGateway(requestedRestPort, rpcAddr)
		go rgw.Serve()

		serverCmdInfo.restServingAddr = rgw.GetRestGatewayAddr()
		serverCmdInfo.restStopFunc = rgw.Close
		serverCmdInfo.restInitialized = true
	}

	serveErr := gs.Serve(helloServer)
	if serveErr != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "error in grpc server Serve(): %v", err)
		return err
	}

	return nil
}

func getRpcServingPort(maxAttempts int) (int, error) {
	serverCmdInfo, err := getServerCmdInfo(maxAttempts)
	if err != nil {
		return -1, err
	}
	return serverCmdInfo.rpcServingPort, nil
}

func getRestServingAddr(maxAttempts int) (net.Addr, error) {
	serverCmdInfo, err := getServerCmdInfo(maxAttempts)
	if err != nil {
		return nil, err
	}
	return serverCmdInfo.restServingAddr, nil
}

func getServerCmdInfo(maxAttempts int) (*ServerCmdInfo, error) {

	for attempts := 1; serverCmdInfo == nil && attempts < maxAttempts; attempts++ {
		time.Sleep(3 * time.Second)
	}

	if serverCmdInfo == nil {
		errMsg := fmt.Sprintf("error not initialized after %d attempts", maxAttempts)
		return nil, errors.New(errMsg)
	}
	return serverCmdInfo, nil

}
