package e2e

import (
	"context"
	"myapp/greeter"
	"myapp/greeterserver"
	"myapp/grpcservice"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/test/bufconn"
)

func TestBufferServing(t *testing.T) {

	assert := assert.New(t)

	buffer := 1024 * 1024
	listener := bufconn.Listen(buffer)

	// first spin up the server
	gs := grpcservice.NewServerListner(listener)
	assert.NotNil(gs)

	helloServer := greeterserver.NewHelloServer()
	defer helloServer.Stop()

	go gs.Serve(helloServer)

	// now let's use the buffer with the client

	cc, bccErr := grpcservice.NewBufferedClientConn(context.Background(), listener)
	assert.NotNil(cc)
	assert.Nil(bccErr)
	verifyClientCalls(t, cc)
}

func TestPortServing(t *testing.T) {

	assert := assert.New(t)

	// first spin up the server

	gs, err := grpcservice.NewServerFromPort(0)
	assert.NotNil(gs)
	assert.Nil(err)
	serverAssignedPort, portErr := gs.GetServicePort()
	assert.Nil(portErr)
	assert.Greater(serverAssignedPort, 0)

	helloServer := greeterserver.NewHelloServer()
	defer helloServer.Stop()
	go gs.Serve(helloServer)

	// now let's use the buffer with the client

	ctx := context.Background()

	cc, bccErr := grpcservice.NewNetClientConn(ctx, "", serverAssignedPort)
	assert.NotNil(cc)
	assert.Nil(bccErr)

	verifyClientCalls(t, cc)
}

func verifyClientCalls(t *testing.T, gccc *grpcservice.Clientconn) {

	assert := assert.New(t)

	var (
		nameInput  string
		timesInput int
		restInput  int

		req greeter.HelloRequest
	)

	// check that the name comes back
	nameInput = "dolly"
	timesInput = 1
	req = greeter.HelloRequest{Name: nameInput, Times: int64(timesInput)}
	replyStream, helloErr := gccc.Call(&req)
	assert.Nil(helloErr)
	assert.NotNil(replyStream)

	replies, err := grpcservice.ReplyStreamToBuffer(replyStream)
	assert.Nil(err)
	assert.Len(replies, int(timesInput))
	assert.Contains(replies[0].GetMessage(), nameInput)

	// check the times parameter
	nameInput = "dolly"
	timesInput = 2
	req = greeter.HelloRequest{Name: nameInput, Times: int64(timesInput)}
	replyStream, err = gccc.Call(&req)
	assert.Nil(err)
	assert.NotNil(replyStream)

	replies, err = grpcservice.ReplyStreamToBuffer(replyStream)
	assert.Nil(err)
	assert.Len(replies, int(timesInput))

	// check the rest parameter
	nameInput = "dolly"
	timesInput = 1
	restInput = 2

	beforeCallTime := time.Now()

	req = greeter.HelloRequest{Name: nameInput, Times: int64(timesInput), Rest: int64(restInput)}
	replyStream, err = gccc.Call(&req)
	assert.Nil(err)
	assert.NotNil(replyStream)

	elapsed := time.Since(beforeCallTime)
	assert.Less(elapsed, time.Duration(restInput)*time.Second)

	replies, err = grpcservice.ReplyStreamToBuffer(replyStream)
	assert.Nil(err)
	assert.Len(replies, int(timesInput))

}
