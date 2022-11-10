package e2e

import (
	"context"
	"io"
	"net"

	pb "myapp/greeter"
	"myapp/greeterserver"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func Greeter(ctx context.Context) (pb.GreeterClient, func()) {
	buffer := 1024 * 1024
	listener := bufconn.Listen(buffer)

	go func() {
		stopFunc, err := greeterserver.ServeListener(listener)
		defer stopFunc()
		if err != nil {
			panic(err)
		}
	}()

	cd := func(context.Context, string) (net.Conn, error) { return listener.Dial() }
	conn, _ := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(cd),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock())

	client := pb.NewGreeterClient(conn)

	return client, func() {} // TODO lots to do
}

func ReplyStreamToBuffer(replyStream pb.Greeter_SayHelloClient) ([]*pb.HelloReply, error) {

	var replies []*pb.HelloReply

	for {
		r, err := replyStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return replies, err
		}
		replies = append(replies, r)
	}

	return replies, nil
}
