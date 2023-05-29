package cmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ServerCommand(t *testing.T) {

	assert := assert.New(t)

	port := 18888 // TODO find an unused port somehow
	portArg := fmt.Sprintf("--port=%d", port)

	go func() {
		cmd := NewServerCmd()
		cmd.SetArgs([]string{portArg})

		GenericCommandRunner(t, cmd)
		assert.Fail("should not reach here, server command should block")
	}()

	clientCmd := NewClientCmd()
	clientCmd.SetArgs([]string{portArg})

	GenericCommandRunner(t, clientCmd, "Hello, World!")

}
