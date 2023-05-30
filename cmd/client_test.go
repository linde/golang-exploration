package cmd

import (
	"fmt"
	"testing"
)

func Test_ClientCommand(t *testing.T) {

	portStr := "--port=8888"

	serverCmd := NewServerCmd()
	serverCmd.SetArgs([]string{portStr})
	go GenericCommandRunner(t, serverCmd)

	// test the port parameter
	portClientCmd := NewClientCmd()
	portClientCmd.SetArgs([]string{portStr})
	GenericCommandRunner(t, portClientCmd, "Hello, world!")

	// TODO test port validation, ie pass in a negative

	// test passing in a host
	hostClientCmd := NewClientCmd()
	hostClientCmd.SetArgs([]string{portStr, "--host=127.0.0.1"})
	GenericCommandRunner(t, hostClientCmd, "Hello, world!")

	// test a specific name
	nameCommandClient := NewClientCmd()
	nameStr := "Sunshine"
	nameCommandClient.SetArgs([]string{portStr, "--name=" + nameStr})
	GenericCommandRunner(t, nameCommandClient, "Hello, "+nameStr)

	// test the --times parameter
	timesCommandClient := NewClientCmd()
	times := 5
	timesCommandClient.SetArgs([]string{portStr, fmt.Sprintf("--times=%d", times)})
	validationStr := fmt.Sprintf("%d of %d", times, times)
	GenericCommandRunner(t, timesCommandClient, validationStr)

	// test the timeout and the wait params

	waitCommandClient := NewClientCmd()
	timeoutSeconds := 1
	timeoutArg := fmt.Sprintf("--timeout=%d", timeoutSeconds)
	waitCommandClient.SetArgs([]string{portStr, "--wait=10", timeoutArg})
	deadlineErrorMsg := fmt.Sprintf("response exceeded deadline of %d seconds", timeoutSeconds)
	GenericCommandRunner(t, waitCommandClient, deadlineErrorMsg)

}
