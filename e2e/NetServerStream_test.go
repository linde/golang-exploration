package e2e

import (
	"context"
	"net"
	"testing"

	pb "myapp/greeter"
	greeterclient "myapp/greeterclient"
	greeterserver "myapp/greeterserver"

	"github.com/stretchr/testify/assert"
)

func TestNetServerStream(t *testing.T) {

	ctx := context.Background()

	assert := assert.New(t)
	assert.NotNil(assert)

	port := 0 // chooses a random open port

	lis, cancel, err := greeterserver.ServePortAsync(port)
	assert.NotNil(cancel)
	// defer cancel()
	assert.Nil(err)

	hostAssignedPort := lis.Addr().(*net.TCPAddr).Port
	assert.Greater(hostAssignedPort, 0)

	cc, _ := greeterclient.NewNetClientConn(ctx, "", hostAssignedPort)
	assert.NotNil(cc)

	var (
		nameInput  string
		timesInput int
		req        pb.HelloRequest
	)

	// check that the name comes back
	nameInput = "dolly"
	timesInput = 1
	req = pb.HelloRequest{Name: nameInput, Times: int64(timesInput)}
	replyStream, helloErr := cc.Call(&req)
	assert.Nil(helloErr)
	assert.NotNil(replyStream)

	replies, err := greeterclient.ReplyStreamToBuffer(replyStream)
	assert.Nil(err)
	assert.Len(replies, int(timesInput))
	assert.Contains(replies[0].GetMessage(), nameInput)

}
