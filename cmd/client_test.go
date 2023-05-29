package cmd

import (
	"testing"
)

func Test_ClientCommand(t *testing.T) {

	portStr := "--port=8888"

	serverCmd := NewServerCmd()
	serverCmd.SetArgs([]string{portStr})
	go GenericCommandRunner(t, serverCmd)

	clientCmd := NewClientCmd()
	clientCmd.SetArgs([]string{portStr})
	GenericCommandRunner(t, clientCmd, "Hello, world!")
}
