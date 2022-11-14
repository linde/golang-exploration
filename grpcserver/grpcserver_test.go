package grpcserver

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/test/bufconn"
)

func TestNetGrpcServer(t *testing.T) {

	assert := assert.New(t)
	assert.NotNil(assert)

	var (
		err error
		gs  *grpcserver
	)
	intendedPort := 10 * 1000 //this is prob going to be flakey
	gs, err = NewServerFromPort(intendedPort)
	assert.NotNil(gs)
	assert.Nil(err)
	serverPort, err := gs.GetServicePort()
	assert.Nil(err)
	assert.NotNil(serverPort)
	assert.Equal(intendedPort, serverPort)

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
