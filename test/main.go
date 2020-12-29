package main

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L. -lposh3
#include <posh3.h>
*/
import "C"
import (
	"bufio"
	"bytes"
	"errors"
	"log"
	"os/exec"
	"strings"
)

type cmd struct{}

func (c *cmd) runCommand(command string) (string, error) {
	cmd := exec.Command(command, "--version")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.Start()
	if err != nil {
		return "", err
	}
	defer func() {
		_ = cmd.Process.Kill()
	}()
	output := new(bytes.Buffer)
	defer output.Reset()
	buf := bufio.NewReader(stdout)
	multiline := false
	for {
		line, _, _ := buf.ReadLine()
		if line == nil {
			break
		}
		if multiline {
			output.WriteString("\n")
		}
		output.Write(line)
		multiline = true
	}
	return output.String(), nil
}

func (c *cmd) runCommandFromRust(command string) (string, error) {
	commandC := C.CString(command)
	response := C.getCommandOutput(commandC)
	defer C.DestroyResponse(response)
	output := C.GoString(response.output)
	err := C.GoString(response.err)
	if err != "" {
		return "", errors.New(err)
	}
	valueClean := strings.TrimSuffix(output, "\n")
	return valueClean, nil
}
