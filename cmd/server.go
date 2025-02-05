package cmd

import (
	"fmt"
	"myapp/grpcservice"
	"myapp/helloserver"
	"myapp/restserver"
	"net"
	"time"

	"github.com/spf13/cobra"
)

// This struct helps us make state available from our normally blocking
// --server command which launche both a grpc server and also optionally
// a rest gateway for it.

type ServerCommand struct {
	requestedRpcPort, servingRpcPort int // TODO make like rest fields
	rpcReady                         bool
	rpcStopFunc                      func()

	restRequestedPort int
	restServingAddr   net.Addr // TODO should this have an accessor with a timeout?
	restReady         bool
	restStopFunc      func()

	Cmd *cobra.Command
}

var serverCmd = NewServerCommand()

func NewServerCommand() *ServerCommand {

	sc := &ServerCommand{}

	sc.Cmd = &cobra.Command{
		Use:   "server",
		Short: "example server for the greeter service",

		// here we inline a RunE function that makes our ServerCommand avail
		RunE: func(*cobra.Command, []string) error {
			return sc.doServerCmd()
		},
	}

	sc.Cmd.Flags().IntVarP(&sc.requestedRpcPort, "port", "p", DEFAULT_PORT, "rpcserver port")
	sc.Cmd.Flags().IntVarP(&sc.restRequestedPort, "rest", "r", -1, "port to use to also enable the rest gateway")

	return sc
}

func init() {
	RootCmd.AddCommand(serverCmd.Cmd)
}

// TODO is this better as a receiver or as a param for a pointer to the struct?
func (sc *ServerCommand) doServerCmd() error {

	// start by creating the grpcservice
	grpcsvc, err := grpcservice.NewServerFromPort(sc.requestedRpcPort)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}
	helloServer := helloserver.NewHelloServer()
	defer helloServer.Stop()

	// collect the serving port used. it might be different
	// if 0 were passed in. also grab the close func
	rpcAddr, err := grpcsvc.GetServiceTCPAddr()
	if err != nil {
		return fmt.Errorf("error determining the serving address: %w", err)
	}
	sc.servingRpcPort = rpcAddr.Port
	sc.rpcStopFunc = helloServer.Stop
	sc.rpcReady = true

	// if restPort is configured, add a rest gateway to our new rpcAddr
	if sc.restRequestedPort >= 0 {
		rgw := restserver.NewRestGateway(sc.restRequestedPort, rpcAddr)
		go rgw.Serve()

		sc.restServingAddr = rgw.GetRestGatewayAddr()
		sc.restStopFunc = rgw.Close
		sc.restReady = true
	}

	// ok, let's go - this will block until closed
	serveErr := grpcsvc.Serve(helloServer)
	if serveErr != nil {
		return fmt.Errorf("error in grpc server Serve(): %w", err)
	}
	return nil
}

func (sc *ServerCommand) GetRunE() (runE func(*cobra.Command, []string) error) {

	return func(*cobra.Command, []string) error {
		return sc.doServerCmd()
	}
}

func (sc *ServerCommand) waitForReady(flag *bool, retries int, retryWait time.Duration) bool {

	for attempts := 0; !*flag && attempts < retries; attempts++ {
		time.Sleep(retryWait)
	}
	return *flag

}

func (sc *ServerCommand) WaitForRpcReady(retries int, retryWait time.Duration) bool {
	return sc.waitForReady(&sc.rpcReady, retries, retryWait)
}

func (sc *ServerCommand) WaitForRestReady(retries int, retryWait time.Duration) bool {
	return sc.waitForReady(&sc.restReady, retries, retryWait)
}

func (sc *ServerCommand) Close() {

	if sc.restStopFunc != nil {
		sc.restStopFunc()
	}
	if sc.rpcStopFunc != nil { // TODO this shouldnt need a guard, right?
		sc.rpcStopFunc()
	}
}
