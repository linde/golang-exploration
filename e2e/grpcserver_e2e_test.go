package e2e

import (
	"context"
	"fmt"
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

func verifyClientCalls(t *testing.T, grpccc *grpcservice.Clientconn) {

	tests := []struct {
		nameInput   string
		timesInput  int64
		restInput   int64
		errExpected bool
	}{
		{"dolly", 1, 0, false},
		{"dolly", 2, 0, false},
		{"dolly", 2, 3, false},
	}

	for _, test := range tests {

		testName := fmt.Sprintf("%s x%d resting %d sec (errExpected=%v)",
			test.nameInput, test.timesInput, test.restInput, test.errExpected)

		t.Run(testName, func(ttt *testing.T) {
			assertNested := assert.New(ttt)

			req := greeter.HelloRequest{Name: test.nameInput, Times: test.timesInput, Rest: test.restInput}

			beforeCallTime := time.Now()

			replyStream, err := grpccc.Call(&req)
			if test.errExpected {
				assertNested.NotNil(err)
				return
			}

			elapsed := time.Since(beforeCallTime)
			if test.restInput > 0 {
				assertNested.Less(elapsed, time.Duration(test.restInput)*time.Second)
			}

			assertNested.Nil(err)
			assertNested.NotNil(replyStream)
			replies, err := grpcservice.ReplyStreamToBuffer(replyStream)
			assertNested.Nil(err)
			assertNested.Len(replies, int(test.timesInput))
			if len(replies) > 0 {
				assertNested.Contains(replies[0].GetMessage(), test.nameInput)
			}
		})
	}

}
