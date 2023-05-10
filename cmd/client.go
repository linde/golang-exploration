package cmd

import (
	"context"
	"fmt"
	"log"
	"myapp/greeter"
	"myapp/grpcservice"
	"time"

	"github.com/spf13/cobra"
)

func NewClientCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "client",
		Short: "minimal grpc client for greeter service",
		Run:   doClientRun,
	}

}

var clientCmd = NewClientCmd()

var name, host string
var times, rest int64
var timeoutSecs int

func init() {
	RootCmd.AddCommand(clientCmd)

	clientCmd.Flags().StringVarP(&host, "host", "s", "localhost", "server host")
	clientCmd.Flags().StringVarP(&name, "name", "n", "world", "whom to greet")
	clientCmd.Flags().Int64VarP(&times, "times", "t", 1, "times to greet them")
	clientCmd.Flags().Int64VarP(&rest, "rest", "r", 0, "seconds to sleep before serving")
	clientCmd.Flags().IntVarP(&timeoutSecs, "timeout", "x", 60, "timeout (in seconds)")
}

func doClientRun(cmd *cobra.Command, args []string) {

	port, _ := cmd.Flags().GetInt("port") // TODO should this be with the others?

	timeout := time.Second * time.Duration(timeoutSecs)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	target := fmt.Sprintf("%s:%d", host, port)
	client, err := grpcservice.NewNetClientConn(ctx, target)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer client.Close()

	request := &greeter.HelloRequest{Name: name, Times: times, Rest: rest}
	replyStream, err := client.Call(request)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	replies, err := grpcservice.ReplyStreamToBuffer(replyStream)
	if err != nil {
		log.Fatalf("could not unbuffer the stream: %v", err)
	}
	for i, reply := range replies {
		log.Printf("got %d: %s", i, reply)
	}

}

/***
func handleReply(r *pb.HelloReply, err error) {
	if err != nil {
		log.Fatalf("client.SayHello failed: %v", err)
	}
	log.Printf("clientCmd: %s", r.GetMessage())
}
****/
