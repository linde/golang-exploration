package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ServerCommand(t *testing.T) {
	assert := assert.New(t)

	cmd := NewServerCmd()
	assert.NotNil(cmd)

	cmd.SetArgs([]string{"--rest=8888"})
	//GenericCommandRunner(t, cmd /*** no assertions bc this wont return ***/)

}
