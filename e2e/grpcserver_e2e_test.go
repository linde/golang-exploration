package e2e

import (
	"context"
	"fmt"
	"math/rand"
	"myapp/greeter"
	"myapp/grpcservice"
	"myapp/helloserver"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func TestBufferServing(t *testing.T) {

	assert := assert.New(t)

	buffer := 1024 * 1024
	listener := bufconn.Listen(buffer)

	// first spin up the server
	gs := grpcservice.NewServerListner(listener)
	assert.NotNil(gs)

	helloServer := helloserver.NewHelloServer()
	defer helloServer.Stop()

	go gs.Serve(helloServer)

	// now let's use the buffer with the client

	cc, bccErr := grpcservice.NewBufferedClientConn(context.Background(), listener)
	assert.NotNil(cc)
	assert.Nil(bccErr)
	verifyClientCalls(t, cc)
}

func TestPortAssignment(t *testing.T) {

	assert := assert.New(t)

	const (
		PORT_MIN, PORT_MAX        int    = 49152, 65535
		ATTEMPTS_TO_GET_FREE_PORT int    = 10
		ADDR_IN_USE_ERROR_MSG     string = "address already in use"
	)

	getRandomFromRange := func(min, max int) int {
		return rand.Intn(max-min) + min
	}

	for attempt := 0; attempt < ATTEMPTS_TO_GET_FREE_PORT; attempt++ {

		// get a port from a range that is reserved for dynamic use. try
		// ATTEMPTS_TO_GET_FREE_PORT times to make sure it's not in use

		assignedPort := getRandomFromRange(PORT_MIN, PORT_MAX)
		gs, err := grpcservice.NewServerFromPort(assignedPort)
		if err != nil && strings.Contains(err.Error(), ADDR_IN_USE_ERROR_MSG) {
			continue
		}
		assert.NotNil(gs)
		assert.Nil(err)

		serverAssignedPort, portErr := gs.GetServicePort()
		assert.Nil(portErr)
		assert.Equal(assignedPort, serverAssignedPort)
		return
	}

	failMsg := fmt.Sprintf("tried %d times, but cannot get free port in range [%d,%d]",
		ATTEMPTS_TO_GET_FREE_PORT, PORT_MIN, PORT_MAX)
	assert.Fail(failMsg)
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

	helloServer := helloserver.NewHelloServer()
	defer helloServer.Stop()
	go gs.Serve(helloServer)

	// now let's use the buffer with the client

	ctx := context.Background()

	target := fmt.Sprintf(":%d", serverAssignedPort)
	cc, bccErr := grpcservice.NewNetClientConn(ctx, target)
	assert.NotNil(cc)
	assert.Nil(bccErr)

	verifyClientCalls(t, cc)
}

func verifyClientCalls(t *testing.T, grpccc *grpcservice.Clientconn) {

	tests := []struct {
		nameInput    string
		timesInput   int64
		pauseInput   int64
		codeExpected codes.Code
	}{
		{"expectingInvalidArgumentError", 1, -1, codes.InvalidArgument},
		{"expectingInvalidArgumentError", -1, 0, codes.InvalidArgument},
		{"dolly", 1, 0, codes.OK},
		{"dolly", 2, 0, codes.OK},
		{"dolly", 2, 3, codes.OK},
	}

	for idx, test := range tests {

		testName := fmt.Sprintf("verifyClientCalls(idx:%d){name:%s,times:%v,rest:%v,expectedCode:%s}",
			idx, test.nameInput, test.timesInput, test.pauseInput, test.codeExpected.String())
		t.Run(testName, func(tt *testing.T) {

			nestedAssert := assert.New(tt)

			req := greeter.HelloRequest{Name: test.nameInput, Times: test.timesInput, Pause: test.pauseInput}

			beforeCallTime := time.Now()

			replyStream, err := grpccc.Call(&req)
			nestedAssert.Nil(err, "unexpected err in grpc Call()")

			elapsed := time.Since(beforeCallTime)
			if test.pauseInput > 0 {
				nestedAssert.Less(elapsed, time.Duration(test.pauseInput)*time.Second)
			}

			nestedAssert.Nil(err)
			nestedAssert.NotNil(replyStream)
			replies, err := grpcservice.ReplyStreamToBuffer(replyStream)

			if test.codeExpected != codes.OK {
				nestedAssert.NotNil(err)
				st := status.Convert(err)
				nestedAssert.Equal(test.codeExpected, st.Code())
				return
			}

			nestedAssert.Nil(err)
			nestedAssert.Len(replies, int(test.timesInput))
			if len(replies) > 0 {
				nestedAssert.Contains(replies[0].GetMessage(), test.nameInput)
			}
		})
	}

}

func BenchmarkBufferServing(b *testing.B) {

	gs, _ := grpcservice.NewServerFromPort(0)
	serverAssignedPort, _ := gs.GetServicePort()
	helloServer := helloserver.NewHelloServer()
	defer helloServer.Stop()
	go gs.Serve(helloServer)

	ctx := context.Background()

	target := fmt.Sprintf(":%d", serverAssignedPort)
	cc, _ := grpcservice.NewNetClientConn(ctx, target)

	for i := 0; i < 1000; i++ {
		testName := fmt.Sprintf("name_%d", i)
		req := greeter.HelloRequest{Name: testName, Times: int64(i), Pause: 0}
		replyStream, _ := cc.Call(&req)

		replies, _ := grpcservice.ReplyStreamToBuffer(replyStream)
		b.Logf("BenchmarkBufferServing got %d replies", len(replies))
	}

}
