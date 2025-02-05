package cmd

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_ClientCommand(t *testing.T) {

	assert := assert.New(t)

	serverCmd := NewServerCommand()
	defer serverCmd.Close()
	serverCmd.Cmd.SetArgs([]string{"--port=0"}) // use zero to grab an open port
	go func() {
		GenericCommandRunner(t, serverCmd.Cmd)
	}()

	rpcReady := serverCmd.WaitForRpcReady(10, 2*time.Second)
	assert.True(rpcReady, "timed out waiting for gRPC service")

	rpcPort := serverCmd.rpcServingPort
	assert.Greater(rpcPort, 0)
	clientPortArg := fmt.Sprintf("--port=%d", rpcPort)

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

	// // test the timeout and the wait params

	waitCommandClient := NewClientCmd()
	timeoutSeconds := 1
	timeoutArg := fmt.Sprintf("--timeout=%d", timeoutSeconds)
	waitCommandClient.SetArgs([]string{clientPortArg, "--wait=10", timeoutArg})
	deadlineErrorMsg := fmt.Sprintf("response exceeded deadline of %d seconds", timeoutSeconds)
	GenericCommandRunner(t, waitCommandClient, deadlineErrorMsg)

}
