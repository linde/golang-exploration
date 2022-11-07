package greeteclient

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "myapp/helloservice"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type greeterclient struct {
	host string
	port int
	c    pb.GreeterClient
}

// TODO initialize the return in the signature. idiomatic no?
func New(host string, port int) *greeterclient {

	gc := greeterclient{host: host, port: port}

	addr := fmt.Sprintf("%s:%d", host, port)

	log.Printf("New() client to: %s", addr)

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	// TODO understand how to do this since the method ends: defer conn.Close()

	gc.c = pb.NewGreeterClient(conn)
	return &gc
}

type ReplyHandler func(reply *pb.HelloReply, err error)

func (gc *greeterclient) Call(name string, times, rest int64, timeoutSecs int, rh ReplyHandler) {

	timeout := time.Second * time.Duration(timeoutSecs)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	request := pb.HelloRequest{Name: name, Times: times, Rest: rest}

	stream, err := gc.c.SayHello(ctx, &request)

	if err != nil {
		log.Fatalf("client.ListFeatures failed: %v", err)
	}

	for {
		r, err := stream.Recv()
		if err == io.EOF {
			break
		}
		rh(r, err)
	}

}
