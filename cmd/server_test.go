package cmd

import (
	"fmt"
	"myapp/testutils"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_ServerCommandRPC(t *testing.T) {

	assert := assert.New(t)

	serverCmd := NewServerCommand()
	defer serverCmd.Close()
	serverCmd.Cmd.SetArgs([]string{"--port=0"}) // use zero to grab an open port
	go func() {
		GenericCommandRunner(t, serverCmd.Cmd)
	}()

	rpcReady := serverCmd.WaitForRpcReady(10, 2*time.Second)
	assert.True(rpcReady, "timed out waiting for gRPC service")

	rpcPort := serverCmd.servingRpcPort
	assert.Greater(rpcPort, 0)
	rpcPortArg := fmt.Sprintf("--port=%d", rpcPort)

	clientCmd := NewClientCmd()
	clientCmd.SetArgs([]string{rpcPortArg})

	GenericCommandRunner(t, clientCmd, "Hello, World!")

}

func Test_ServerCommandRestGateway(t *testing.T) {

	assert := assert.New(t)

	serverCmd := NewServerCommand()
	defer serverCmd.Close()

	// in both cases, use zero to grab open ports
	serverCmd.Cmd.SetArgs([]string{"--port=0", "--rest=0"})

	go func() {
		GenericCommandRunner(t, serverCmd.Cmd)
	}()

	restReady := serverCmd.WaitForRestReady(10, 2*time.Second)
	assert.True(restReady, "Took too long for the rest gateway to become available")

	restAddr := serverCmd.restServingAddr
	assert.NotNil(restAddr)

	nameInput := "Cornelius"
	url := fmt.Sprintf("http://%s/v1/helloservice/sayhello?name=%s&times=1", restAddr.String(), nameInput)

	testutils.DoHttpTest(t, url, http.StatusOK, []string{nameInput})

}
