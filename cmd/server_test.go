package cmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// make this true if you want to run the test below to debug interactively
const RUN_BLOCKING_TEST = false

func Test_ServerCommand(t *testing.T) {
	assert := assert.New(t)

	if RUN_BLOCKING_TEST != true {
		return
	}

	cmd := NewServerCmd()
	assert.NotNil(cmd)

	port := 18888
	portArg := fmt.Sprintf("--port=%d", port)
	cmd.SetArgs([]string{portArg})
	go func() {
		out := GenericCommandRunner(t, cmd /*** no assertions bc this wont return ***/)
		t.Errorf("got: %s", out)
	}()

	clientCmd := NewClientCmd()
	clientCmd.SetArgs([]string{portArg})

	GenericCommandRunner(t, clientCmd)

}
