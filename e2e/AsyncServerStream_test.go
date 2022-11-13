package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/test/bufconn"

	pb "myapp/greeter"

	greeterclient "myapp/greeterclient"
	greeterserver "myapp/greeterserver"
)

func TestAsyncServerStream(t *testing.T) {

	ctx := context.Background()
	assert := assert.New(t)

	buffer := 1024 * 1024
	listener := bufconn.Listen(buffer)

	// run a buffered Server
	cancelFunc := greeterserver.ServeListenerAsync(listener)
	assert.NotNil(cancelFunc)
	// TODO why can't i do this: defer cancelFunc()

	bufclientConn, bccErr := greeterclient.NewBufferedClientConn(ctx, listener)
	assert.NotNil(bufclientConn)
	assert.Nil(bccErr)

	var (
		nameInput  string
		timesInput int
		restInput  int64

		req pb.HelloRequest
	)

	// check that the name comes back
	nameInput = "dolly"
	timesInput = 1
	req = pb.HelloRequest{Name: nameInput, Times: int64(timesInput)}
	replyStream, helloErr := bufclientConn.Call(&req)
	assert.Nil(helloErr)
	assert.NotNil(replyStream)

	replies, err := greeterclient.ReplyStreamToBuffer(replyStream)
	assert.Nil(err)
	assert.Len(replies, int(timesInput))
	assert.Contains(replies[0].GetMessage(), nameInput)

	// check the times parameter
	nameInput = "dolly"
	timesInput = 2
	req = pb.HelloRequest{Name: nameInput, Times: int64(timesInput)}
	replyStream, err = bufclientConn.Call(&req)
	assert.Nil(err)
	assert.NotNil(replyStream)

	replies, err = greeterclient.ReplyStreamToBuffer(replyStream)
	assert.Nil(err)
	assert.Len(replies, int(timesInput))

	// check the rest parameter
	nameInput = "dolly"
	timesInput = 1
	restInput = 2

	beforeCallTime := time.Now()

	req = pb.HelloRequest{Name: nameInput, Times: int64(timesInput), Rest: restInput}
	replyStream, err = bufclientConn.Call(&req)
	assert.Nil(err)
	assert.NotNil(replyStream)

	elapsed := time.Since(beforeCallTime)
	assert.Less(elapsed, time.Duration(restInput)*time.Second)

	replies, err = greeterclient.ReplyStreamToBuffer(replyStream)
	assert.Nil(err)
	assert.Len(replies, int(timesInput))
}
