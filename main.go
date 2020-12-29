package command

/*
#cgo LDFLAGS: -L./lib -lcommand
#include <stdlib.h>
#include "./lib/command.h"
*/
import "C"
import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
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
	value := C.GoString(response.output)
	const err string = "err: "
	if strings.HasPrefix(value, err) {
		errStr := strings.TrimPrefix(value, err)
		return "", errors.New(errStr)
	}
	valueClean := strings.TrimSuffix(value, "\n")
	return valueClean, nil
}

func main() {
	cmd := &cmd{}
	command := "git"
	start := time.Now()
	output, err := cmd.runCommandFromRust(command)
	elapsed := time.Since(start)
	log.Printf("runCommandFromRust took %s", elapsed)
	start = time.Now()
	output, err = cmd.runCommand(command)
	elapsed = time.Since(start)
	log.Printf("runCommand took %s", elapsed)
	if err != nil {
		fmt.Println("Found error: ", err)
		return
	}
	fmt.Println(output)
}
