package greeterserver

import (
	"context"
	"fmt"
	"log"
	pb "myapp/helloservice"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v, %d times", in.GetName(), in.GetTimes())
	returnMessage := fmt.Sprintf("Hiya %s %d times", in.GetName(), in.GetTimes())
	return &pb.HelloReply{Message: returnMessage}, nil
}

func New() *server {
	return &server{}
}
