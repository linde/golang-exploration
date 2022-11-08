package greeterserver

import (
	"fmt"
	"log"
	pb "myapp/helloservice"
	"time"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(in *pb.HelloRequest, stream pb.Greeter_SayHelloServer) error {
	log.Printf("Received: %v, %d times, after %d seconds rest", in.GetName(), in.GetTimes(), in.GetRest())

	time.Sleep(time.Duration(in.GetRest()) * time.Second)

	if in.GetTimes() > 0 {
		for i := int64(0); i < in.GetTimes(); i++ {

			returnMessage := fmt.Sprintf("Hello, %s! (%d of %d)", in.GetName(), i+1, in.GetTimes())
			if err := stream.Send(&pb.HelloReply{Message: returnMessage}); err != nil {
				return err
			}
		}
	}
	return nil
}

func New() *server {
	return &server{}
}
