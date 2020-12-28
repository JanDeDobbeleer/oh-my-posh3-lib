package main

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
	"unsafe"
)

func runCommand(command string) (string, error) {
	cmd := exec.Command(command, "Hello world")
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

func runCommandFromRust(command string) (string, error) {
	commandC := C.CString(command)
	defer C.free(unsafe.Pointer(commandC))
	response := C.getCommandOutput(commandC)
	defer C.free(unsafe.Pointer(response))
	value := C.GoString(response)
	const err string = "err: "
	if strings.HasPrefix(value, err) {
		errStr := strings.TrimPrefix(value, err)
		return "", errors.New(errStr)
	}
	return value, nil
}

func main() {
	command := "echo"
	start := time.Now()
	output, err := runCommandFromRust(command)
	elapsed := time.Since(start)
	log.Printf("runCommandFromRust took %s", elapsed)
	start = time.Now()
	output, err = runCommand(command)
	elapsed = time.Since(start)
	log.Printf("runCommand took %s", elapsed)
	if err != nil {
		fmt.Println("Found error: ", err)
		return
	}
	fmt.Println(output)
}
