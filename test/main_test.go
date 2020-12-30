package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandOneArg(t *testing.T) {
	cmd := &cmd{}
	command := "git"
	args := []string{"--version"}
	outputGo, errGo := cmd.runCommand(command, args...)
	output, err := cmd.runCommandFromRust(command, args...)
	assert.Equal(t, outputGo, output)
	assert.Equal(t, errGo, err)
}

func TestCommandMultipleArgs(t *testing.T) {
	cmd := &cmd{}
	command := "git"
	args := []string{"status", "-unormal", "--short", "--branch"}
	outputGo, errGo := cmd.runCommand(command, args...)
	output, err := cmd.runCommandFromRust(command, args...)
	assert.Equal(t, outputGo, output)
	assert.Equal(t, errGo, err)
}

func TestCommandNoArgs(t *testing.T) {
	cmd := &cmd{}
	command := "git"
	outputGo, errGo := cmd.runCommand(command)
	output, err := cmd.runCommandFromRust(command)
	assert.Equal(t, outputGo, output)
	assert.Equal(t, errGo, err)
}

func TestCommandError(t *testing.T) {
	cmd := &cmd{}
	command := "git"
	args := []string{"burp"}
	outputGo, errGo := cmd.runCommand(command, args...)
	output, err := cmd.runCommandFromRust(command, args...)
	assert.Equal(t, outputGo, output)
	assert.Equal(t, errGo, err)
}
