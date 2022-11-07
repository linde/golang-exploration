package greeterserver

import (
	"fmt"
	"log"
	pb "myapp/helloservice"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(in *pb.HelloRequest, stream pb.Greeter_SayHelloServer) error {
	log.Printf("Received: %v, %d times", in.GetName(), in.GetTimes())

	// TODO: fix this logic so it skips first reply if less than 1 time, ie --times=0
	for i := int64(0); i < in.GetTimes(); i++ {

		returnMessage := fmt.Sprintf("Hello, %s! (%d of %d)", in.GetName(), i+1, in.GetTimes())
		if err := stream.Send(&pb.HelloReply{Message: returnMessage}); err != nil {
			return err
		}
	}
	return nil
}

func New() *server {
	return &server{}
}
