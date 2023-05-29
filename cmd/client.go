package cmd

import (
	"context"
	"fmt"
	"myapp/greeter"
	"myapp/grpcservice"
	"time"

	"github.com/spf13/cobra"
)

var clientCmd = NewClientCmd()

var name, host string
var times, pause int64
var timeoutSecs int

func NewClientCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "client",
		Short: "minimal grpc client for greeter service",
		RunE:  doClientRun,
	}

	cmd.Flags().StringVarP(&host, "host", "s", "localhost", "server host")
	cmd.Flags().IntVarP(&rpcPort, "port", "p", DEFAULT_PORT, "rpcserver port")
	cmd.Flags().StringVarP(&name, "name", "n", "world", "whom to greet")
	cmd.Flags().Int64VarP(&times, "times", "t", 1, "times to greet them")
	cmd.Flags().Int64VarP(&pause, "wait", "w", 0, "seconds to wait before serving")
	cmd.Flags().IntVarP(&timeoutSecs, "timeout", "x", 60, "timeout (in seconds)")

	return cmd
}

func init() {
	RootCmd.AddCommand(clientCmd)
}

func doClientRun(cmd *cobra.Command, args []string) error {

	timeout := time.Second * time.Duration(timeoutSecs)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	target := fmt.Sprintf("%s:%d", host, rpcPort)
	client, err := grpcservice.NewNetClientConn(ctx, target)
	if err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "did not connect: %v", err)
		return err
	}
	defer client.Close()

	request := &greeter.HelloRequest{Name: name, Times: times, Pause: pause}
	replyStream, err := client.Call(request)
	if err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "did not connect: %v", err)
		return err
	}

	replies, err := grpcservice.ReplyStreamToBuffer(replyStream)
	if err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "could not unbuffer the stream: %v", err)
		return err
	}

	for i, reply := range replies {
		fmt.Fprintf(cmd.OutOrStdout(), "got %d: %s\n", i, reply)
	}

	return nil
}
