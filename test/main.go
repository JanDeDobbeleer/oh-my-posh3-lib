package main

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L. -lposh3
#include <stdlib.h>
#include <posh3.h>
*/
import "C"
import (
	"bufio"
	"bytes"
	"io"
	"os/exec"
	"strings"
	"unsafe"
)

type commandError struct {
	err      string
	exitCode int
}

func (e *commandError) Error() string {
	return e.err
}

type cmd struct{}

func (c *cmd) runCommand(command string, args ...string) (string, error) {
	getOutputString := func(io io.ReadCloser) string {
		output := new(bytes.Buffer)
		defer output.Reset()
		buf := bufio.NewReader(io)
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
		return output.String()
	}
	cmd := exec.Command(command, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", &commandError{
			err:      err.Error(),
			exitCode: 666,
		}
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", &commandError{
			err:      err.Error(),
			exitCode: 667,
		}
	}
	err = cmd.Start()
	if err != nil {
		return "", &commandError{
			err:      err.Error(),
			exitCode: 668,
		}
	}
	defer func() {
		_ = cmd.Process.Kill()
	}()
	stdoutString := getOutputString(stdout)
	stderrString := getOutputString(stderr)
	if stderrString != "" {
		// only wait in case of error reduces the lead time on successful
		// commands on windows due to not calling process.Wait()
		_ = cmd.Wait()
		return "", &commandError{
			err:      stderrString,
			exitCode: cmd.ProcessState.ExitCode(),
		}
	}
	return stdoutString, nil
}

func (c *cmd) runCommandFromRust(command string, args ...string) (string, error) {
	cleanOutput := func(output string) string {
		return strings.TrimSuffix(output, "\n")
	}
	commandC := C.CString(command)
	defer C.free(unsafe.Pointer(commandC))
	var argsC *C.char
	if args != nil {
		argsJoined := strings.Join(args, ";")
		argsC = C.CString(argsJoined)
		defer C.free(unsafe.Pointer(argsC))
	}
	response := C.getCommandOutput(commandC, argsC)
	defer C.DestroyResponse(response)
	output := C.GoString(response.output)
	err := C.GoString(response.err)
	if err != "" {
		return "", &commandError{
			err:      cleanOutput(err),
			exitCode: int(response.status_code),
		}
	}
	return cleanOutput(output), nil
}
