package grpcserver

import (
	"context"
	"myapp/greeterclient"
	"myapp/greeterserver"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/test/bufconn"

	pb "myapp/greeter"
)

func TestNetGrpcServer(t *testing.T) {

	assert := assert.New(t)
	assert.NotNil(assert)

	var (
		err error
		gs  *grpcserver
	)

	// TODO how to test specifying a port without picking a port in use?

	openPort := 0 // random open port
	gs, err = NewServerFromPort(openPort)
	assert.NotNil(gs)
	assert.Nil(err)
	serverAssignedPort, err := gs.GetServicePort()
	assert.Nil(err)
	assert.Positive(serverAssignedPort)
	assert.NotEqual(openPort, serverAssignedPort)

}

func TestBufferGrpcServer(t *testing.T) {

	assert := assert.New(t)
	assert.NotNil(assert)

	var (
		err error
		gs  *grpcserver
	)

	buffer := 1024 * 1024
	listener := bufconn.Listen(buffer)

	gs = NewServerListner(listener)
	assert.NotNil(gs)

	// check to make sure we're throwing a error for a port check for a buff listner
	invalidPort, err := gs.GetServicePort()
	assert.NotNil(err)
	assert.Contains(err.Error(), "expected net.TCPAddr")
	assert.NotNil(invalidPort)
	assert.Less(invalidPort, 0)

}

func TestBufferServing(t *testing.T) {

	ctx := context.Background()
	assert := assert.New(t)

	buffer := 1024 * 1024
	listener := bufconn.Listen(buffer)

	// first spin up the server
	gs := NewServerListner(listener)
	assert.NotNil(gs)

	helloServer := greeterserver.NewHelloServer()
	defer helloServer.Stop()

	go gs.Serve(helloServer)

	// now let's use the buffer with the client
	bufclientConn, bccErr := greeterclient.NewBufferedClientConn(ctx, listener)
	assert.NotNil(bufclientConn)
	assert.Nil(bccErr)

	var (
		nameInput  string
		timesInput int

		req pb.HelloRequest
	)

	// check that the name comes back
	nameInput = "dolly"
	timesInput = 1
	req = pb.HelloRequest{Name: nameInput, Times: int64(timesInput)}
	replyStream, helloErr := bufclientConn.Call(&req)
	assert.Nil(helloErr)
	assert.NotNil(replyStream)

}
