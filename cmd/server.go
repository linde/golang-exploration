package cmd

import (
	"fmt"
	"log"
	"net"

	greeterserver "myapp/greeterserver"
	pb "myapp/helloservice"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "minimal grpc client for greeter service",
	Run:   doServerRun,
}

func init() {
	RootCmd.AddCommand(serverCmd)
}

func doServerRun(cmd *cobra.Command, args []string) {

	portParam, _ := cmd.Flags().GetInt("port")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", portParam))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	server := greeterserver.New()

	pb.RegisterGreeterServer(s, server)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
