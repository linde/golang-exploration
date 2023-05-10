package grpcservice

import (
	"context"
	"io"
	"log"
	"net"

	"myapp/greeter"

	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Clientconn struct {
	conn *grpc.ClientConn
	ctx  context.Context
}

func (cc Clientconn) GetClientConn() *grpc.ClientConn {
	return cc.conn
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

func NewNetClientConn(ctx context.Context, target string) (netcc *Clientconn, returnErr error) {

	log.Printf("NewNetClientConn() client to: %s", target)

	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
		st := status.Convert(err)
		log.Printf("greeterclient.Call had error, status: %v", st)
		stream = nil
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
