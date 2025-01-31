package cmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ClientCommand(t *testing.T) {

	assert := assert.New(t)

	serverCmd := NewServerCmd()
	defer ServerCommandClose()
	serverCmd.SetArgs([]string{"--port=0"}) // port 0 finds an unused port
	go GenericCommandRunner(t, serverCmd)
	defer ServerCommandClose()

	clientPort, clientPortErr := getRpcServingPort(10)
	assert.Nil(clientPortErr)
	assert.Greater(clientPort, 0)
	clientPortArg := fmt.Sprintf("--port=%d", clientPort)

	// test the port parameter
	portClientCmd := NewClientCmd()
	portClientCmd.SetArgs([]string{clientPortArg})
	GenericCommandRunner(t, portClientCmd, "Hello, world!")

	// TODO test port validation, ie pass in a negative

	// test passing in a host
	hostClientCmd := NewClientCmd()
	hostClientCmd.SetArgs([]string{clientPortArg, "--host=127.0.0.1"})
	GenericCommandRunner(t, hostClientCmd, "Hello, world!")

	// test a specific name
	nameCommandClient := NewClientCmd()
	nameStr := "Sunshine"
	nameCommandClient.SetArgs([]string{clientPortArg, "--name=" + nameStr})
	GenericCommandRunner(t, nameCommandClient, "Hello, "+nameStr)

	// test the --times parameter
	timesCommandClient := NewClientCmd()
	times := 5
	timesCommandClient.SetArgs([]string{clientPortArg, fmt.Sprintf("--times=%d", times)})
	validationStr := fmt.Sprintf("%d of %d", times, times)
	GenericCommandRunner(t, timesCommandClient, validationStr)

	// test the timeout and the wait params

	waitCommandClient := NewClientCmd()
	timeoutSeconds := 1
	timeoutArg := fmt.Sprintf("--timeout=%d", timeoutSeconds)
	waitCommandClient.SetArgs([]string{clientPortArg, "--wait=10", timeoutArg})
	deadlineErrorMsg := fmt.Sprintf("response exceeded deadline of %d seconds", timeoutSeconds)
	GenericCommandRunner(t, waitCommandClient, deadlineErrorMsg)

}
