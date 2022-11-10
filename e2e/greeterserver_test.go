package e2e

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	pb "myapp/greeter"
)

func TestGreeterGrpc(t *testing.T) {

	ctx := context.Background()
	assert := assert.New(t)

	client, closer := Greeter(ctx)
	defer closer()

	// check that the name comes back
	var nameInput = "dolly"
	var timesInput = 1
	var req = pb.HelloRequest{Name: nameInput, Times: int64(timesInput), Rest: 0}
	var replyStream, err = client.SayHello(ctx, &req)
	assert.Nil(err)
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
