package greeterserver

import (
	"fmt"
	"log"
	pb "myapp/helloservice"
	"net"
	"time"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func ServePort(port int) error {

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Printf("greeterserver.ServePort() failed to listen: %v", err)
		return err
	}
	return ServeListener(lis)
}

func ServeListener(lis net.Listener) error {

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Printf("greeterserver.ServeListener() failed to listen: %v", err)
		return err
	}
	return nil
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
