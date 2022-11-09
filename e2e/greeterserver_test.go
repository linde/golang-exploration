package e2e

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	gs "myapp/greeterserver"
	pb "myapp/helloservice"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func TestGreeterGrpc(t *testing.T) {

	ctx := context.Background()
	assert := assert.New(t)

	client, closer := greeter(ctx)
	defer closer()

	nameInput := "dolly"
	replyStream, err := client.SayHello(ctx, &pb.HelloRequest{Name: nameInput, Times: 1, Rest: 0})
	assert.Nil(err)
	assert.NotNil(replyStream)

	r, err := replyStream.Recv()
	assert.Nil(err)
	assert.Contains(r.GetMessage(), nameInput)
	assert.Nil(replyStream.Recv())

}

func greeter(ctx context.Context) (pb.GreeterClient, func()) {
	buffer := 1024 * 1024
	listener := bufconn.Listen(buffer)

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &gs.Server{})
	go func() {
		if err := s.Serve(listener); err != nil {
			panic(err)
		}
	}()

	cd := func(context.Context, string) (net.Conn, error) { return listener.Dial() }
	conn, _ := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(cd),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock())

	client := pb.NewGreeterClient(conn)

	return client, s.Stop
}
