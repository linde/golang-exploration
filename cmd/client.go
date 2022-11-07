package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "myapp/helloservice"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "minimal grpc client for greeter service",
	Run:   doClientRun,
}

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

	portParam, _ := cmd.Flags().GetInt("port")
	addr := fmt.Sprintf("%s:%d", host, portParam)

	log.Printf("Greeting %s %d times @ %s after %d seconds rest\n", name, times, addr, rest)

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	timeout := time.Second * time.Duration(timeoutSecs)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	request := pb.HelloRequest{Name: name, Times: times, Rest: rest}

	stream, err := c.SayHello(ctx, &request)
	if err != nil {
		log.Fatalf("client.ListFeatures failed: %v", err)
	}

	for {
		r, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("client.SayHello failed: %v", err)
		}
		log.Printf("Greeting: %s", r.GetMessage())
	}

}
