package cmd

import (
	"context"
	"fmt"
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
var times int64

func init() {
	RootCmd.AddCommand(clientCmd)

	clientCmd.Flags().StringVarP(&host, "host", "s", "localhost", "server host")
	clientCmd.Flags().StringVarP(&name, "name", "n", "world", "whom to greet")
	clientCmd.Flags().Int64VarP(&times, "times", "t", 1, "times to greet them")
}

func doClientRun(cmd *cobra.Command, args []string) {

	portParam, _ := cmd.Flags().GetInt("port")

	addr := fmt.Sprintf("%s:%d", host, portParam)

	fmt.Printf("Greeting %s %d times @ %s\n", name, times, addr)

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name, Times: times})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())

}
