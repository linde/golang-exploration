package greeterclient

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "myapp/greeter"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type greeterclient struct {
	host string
	port int
	conn *grpc.ClientConn
}

// TODO initialize the return in the signature. idiomatic no?
func New(host string, port int) (*greeterclient, error) {

	gc := greeterclient{host: host, port: port}

	addr := fmt.Sprintf("%s:%d", host, port)

	log.Printf("New() client to: %s", addr)

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	gc.conn = conn
	return &gc, err
}

type ReplyHandler func(reply *pb.HelloReply, err error)

func (gc *greeterclient) Call(name string, times, rest int64, timeoutSecs int, rh ReplyHandler) {

	ngc := pb.NewGreeterClient(gc.conn)

	timeout := time.Second * time.Duration(timeoutSecs)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	request := pb.HelloRequest{Name: name, Times: times, Rest: rest}

	stream, err := ngc.SayHello(ctx, &request)

	if err != nil {
		log.Fatalf("greeterclient.Call failed: %v", err)
	}

	for {
		r, err := stream.Recv()
		if err == io.EOF {
			break
		}
		rh(r, err)
	}

}

func (gc *greeterclient) Close() {
	if gc.conn != nil {
		gc.conn.Close()
	}
}
