package cmd

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_RootCommand(t *testing.T) {

	helpAssertionString := "cli"

	cmd := NewRootCmd()
	GenericCommandRunner(t, cmd, helpAssertionString)

	cmd.SetArgs([]string{"--help"})
	GenericCommandRunner(t, cmd, helpAssertionString)

	falseFlag := "--not-a-real-arg"
	cmd.SetArgs([]string{falseFlag})
	GenericCommandRunner(t, cmd, "unknown flag: "+falseFlag)

}

func GenericCommandRunner(t *testing.T, cmd *cobra.Command, outputAssertions ...string) string {
	assert := assert.New(t)

	assert.NotNil(cmd)

	b := bytes.NewBufferString("")

	cmd.SetOut(b)
	cmd.SetErr(b)
	cmd.Execute()
	out, _ := io.ReadAll(b)

	for _, oa := range outputAssertions {
		assert.Contains(strings.ToLower(string(out)), strings.ToLower(oa))
	}
	return string(out)
}
