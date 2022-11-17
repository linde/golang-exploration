package grpcservice

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"myapp/greeter"

	"google.golang.org/grpc/test/bufconn"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Clientconn struct {
	conn *grpc.ClientConn
	ctx  context.Context
}

func NewBufferedClientConn(ctx context.Context, listener *bufconn.Listener) (netcc *Clientconn, returnErr error) {

	cd := func(context.Context, string) (net.Conn, error) { return listener.Dial() }
	conn, _ := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(cd),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock())

	netcc = &Clientconn{conn: conn, ctx: ctx}
	return netcc, nil
}

func NewNetClientConn(ctx context.Context, host string, port int) (netcc *Clientconn, returnErr error) {

	addr := fmt.Sprintf("%s:%d", host, port)

	log.Printf("NewNetClientConn() client to: %s", addr)

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("NewNetClientConn() error: %v", err)
		return nil, err
	}

	netcc = &Clientconn{conn: conn, ctx: ctx}
	return netcc, nil
}

func (gc *Clientconn) Call(req *greeter.HelloRequest) (greeter.Greeter_SayHelloClient, error) {

	ngc := greeter.NewGreeterClient(gc.conn)
	stream, err := ngc.SayHello(gc.ctx, req)
	if err != nil {
		log.Fatalf("greeterclient.Call failed: %v", err)
	}

	return stream, err
}

func ReplyStreamToBuffer(replyStream greeter.Greeter_SayHelloClient) ([]*greeter.HelloReply, error) {

	var replies []*greeter.HelloReply

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

func (gc *Clientconn) Close() {
	if gc.conn != nil {
		gc.conn.Close()
	}
}
