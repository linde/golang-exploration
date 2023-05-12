// TODO rename this package when the server is moved to grpc
package greeterserver

import (
	"fmt"
	"log"
	pb "myapp/greeter"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func NewHelloServer() *grpc.Server {

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	reflection.Register(s)

	return s
}

func (s *server) SayHello(in *pb.HelloRequest, stream pb.Greeter_SayHelloServer) error {

	name, times, pause := in.GetName(), in.GetTimes(), in.GetPause()
	log.Printf("SayHello: %v, %d times, after %d seconds pause", name, times, pause)

	if pause < 0 {
		st := status.Newf(codes.InvalidArgument, "got negative rest duration (%v)", pause)
		return st.Err()
	}
	if times < 0 {
		st := status.Newf(codes.InvalidArgument, "got negative times (%v)", times)
		return st.Err()
	}
	time.Sleep(time.Duration(pause) * time.Second)

	if times > 0 {
		for i := int64(0); i < times; i++ {

			returnMessage := fmt.Sprintf("Hello, %s! (%d of %d)", name, i+1, times)
			if err := stream.Send(&pb.HelloReply{Message: returnMessage}); err != nil {
				return err
			}
		}
	}
	return nil
}
