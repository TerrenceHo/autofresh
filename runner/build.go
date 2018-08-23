package runner

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

// Build runs a command to build a program. First word in buildCommand must be a
// executable command, the others are arguments. Writes output and errors to the
// console. Waits for command to finish before exiting.
func Build(buildString string) (string, error) {
	var cmd *exec.Cmd
	buildArr := strings.Split(buildString, " ")
	buildCmd := buildArr[0]
	buildArgs := buildArr[1:]
	if len(buildArr) > 1 {
		cmd = exec.Command(buildCmd, buildArgs...)
	} else {
		cmd = exec.Command(buildCmd)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Error with StdoutPipe, error: %s\n", err.Error())
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatalf("Error with StderrPipe, error: %s\n", err.Error())
	}

	err = cmd.Start()
	if err != nil {
		log.Fatalf("Error with starting Command, error: %s\n", err.Error())
	}

	io.Copy(os.Stdout, stdout)            // copy output from build command to console
	buildErr, _ := ioutil.ReadAll(stderr) // read all errors into buildErr

	if err := cmd.Wait(); err != nil {
		return string(buildErr), err
	}

	return "", nil
}
