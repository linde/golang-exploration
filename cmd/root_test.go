package cmd

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ExecuteCommand(t *testing.T) {
	assert := assert.New(t)

	cmd := NewRootCmd()
	assert.NotNil(cmd)

	b := bytes.NewBufferString("")

	cmd.SetOut(b)
	cmd.Execute()
	out, err := ioutil.ReadAll(b)
	assert.Nil(err)
	assert.Contains(string(out), "CLI")

}
