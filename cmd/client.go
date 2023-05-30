package cmd

import (
	"context"
	"fmt"
	"myapp/greeter"
	"myapp/grpcservice"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	// TODO rename the proto for sleep to match this param
	cmd.Flags().Int64VarP(&pause, "wait", "w", 0, "seconds to wait before serving")
	cmd.Flags().IntVarP(&timeoutSecs, "timeout", "x", 0, "timeout (in seconds)")

	return cmd
}

func init() {
	RootCmd.AddCommand(clientCmd)
}

func doClientRun(cmd *cobra.Command, args []string) error {

	var ctx = context.Background()
	if timeoutSecs > 0 {
		timeout := time.Second * time.Duration(timeoutSecs)
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
		defer cancel()
	}

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

		// first check if we had a timeoout and exceeded it and report appropriately
		st := status.Convert(err)
		if timeoutSecs > 0 && st.Code() == codes.DeadlineExceeded {
			fmt.Fprintf(cmd.ErrOrStderr(), "Response exceeded deadline of %d seconds\n", timeoutSecs)
			return nil
		}

		fmt.Fprintf(cmd.ErrOrStderr(), "could not unbuffer the stream: %v", err)
		return err
	}

	for i, reply := range replies {
		fmt.Fprintf(cmd.OutOrStdout(), "got %d: %s\n", i, reply)
	}

	return nil
}
