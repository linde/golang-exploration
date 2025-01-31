package cmd

import (
	"fmt"
	"myapp/testutils"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ServerCommandRPC(t *testing.T) {

	assert := assert.New(t)

	serverCmd := NewServerCmd()
	port := 0 // use zero to grab an open port
	portArg := fmt.Sprintf("--port=%d", port)
	serverCmd.SetArgs([]string{portArg})
	defer ServerCommandClose()

	go func() {
		GenericCommandRunner(t, serverCmd)
	}()

	clientPort, clientPortErr := getRpcServingPort(10)
	assert.Nil(clientPortErr)
	assert.Greater(clientPort, 0)
	clientPortArg := fmt.Sprintf("--port=%d", clientPort)

	clientCmd := NewClientCmd()
	clientCmd.SetArgs([]string{clientPortArg})

	GenericCommandRunner(t, clientCmd, "Hello, World!")

}

func Test_ServerCommandRestGateway(t *testing.T) {

	assert := assert.New(t)

	// 0 will find any open port so we use it for both
	portArg := fmt.Sprintf("--port=%d", 0)
	restPortArg := fmt.Sprintf("--rest=%d", 0)
	cmd := NewServerCmd()
	cmd.SetArgs([]string{portArg, restPortArg})

	go GenericCommandRunner(t, cmd)

	restAddr, clientAddrErr := getRestServingAddr(10)
	assert.Nil(clientAddrErr)

	nameInput := "Cornelius"
	url := fmt.Sprintf("http://%s/v1/helloservice/sayhello?name=%s&times=1", restAddr.String(), nameInput)

	testutils.DoHttpTest(t, url, http.StatusOK, []string{nameInput})

}
