package command

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCommandFromRust(t *testing.T) {
	cmd := &cmd{}
	command := "git"
	start := time.Now()
	output, err := cmd.runCommandFromRust(command)
	elapsedRust := time.Since(start)
	start = time.Now()
	outputGo, errGo := cmd.runCommand(command)
	elapsedGo := time.Since(start)
	assert.Equal(t, outputGo, output)
	assert.Equal(t, errGo, err)
	assert.LessOrEqual(t, elapsedGo.Microseconds(), elapsedRust.Microseconds())
}
