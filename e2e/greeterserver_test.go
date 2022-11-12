package e2e

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/test/bufconn"

	pb "myapp/greeter"

	greeterclient "myapp/greeterclient"
	greeterserver "myapp/greeterserver"
)

func TestGreeterGrpc(t *testing.T) {

	ctx := context.Background()
	assert := assert.New(t)

	buffer := 1024 * 1024
	listener := bufconn.Listen(buffer)

	client, gcErr := greeterclient.NewBufferGreeterClient(listener)
	assert.NotNil(gcErr)

	// check that the name comes back
	var nameInput = "dolly"
	var timesInput = 1
	var req = pb.HelloRequest{Name: nameInput, Times: int64(timesInput), Rest: 0}

	var replyBuff = greeterclient.NewReplyStreamBuffer()
	replyStream, helloErr := client.CallRequest(&req, 0, replyBuff.GetReplyHandler())
	assert.Nil(helloErr)
	assert.NotNil(replyStream)

	replies, err := ReplyStreamToBuffer(replyStream)
	assert.Nil(err)
	assert.Len(replies, int(timesInput))
	assert.Contains(replies[0].GetMessage(), nameInput)

	// check the times parameter
	timesInput = 2
	req = pb.HelloRequest{Name: nameInput, Times: int64(timesInput), Rest: 0}
	replyStream, err = client.SayHello(ctx, &req)
	assert.Nil(err)
	assert.NotNil(replyStream)

	replies, err = ReplyStreamToBuffer(replyStream)
	assert.Nil(err)
	assert.Len(replies, int(timesInput))
}

// func ServeListenerAsync(cancel chan (bool), lis net.Listener) error {

func TestAsync(t *testing.T) {

	assert := assert.New(t)
	assert.NotNil(assert)

	buffer := 1024 * 1024
	listener := bufconn.Listen(buffer)

	cancelFunc := greeterserver.ServeListenerAsync(listener)
	defer cancelFunc()

	ctx := context.Background()
	client := GreeterClient(ctx, listener)

	// check the times parameter
	var timesInput = 2
	var req = pb.HelloRequest{Name: "dolly", Times: int64(timesInput), Rest: 0}
	var replyStream, err = client.SayHello(ctx, &req)
	assert.Nil(err)
	assert.NotNil(replyStream)

	var replies, err2 = ReplyStreamToBuffer(replyStream)
	assert.Nil(err2)
	assert.Len(replies, int(timesInput))
}
