package runner

import (
	"io"
	"log"
	"os/exec"

	"github.com/TerrenceHo/autofresh/logger"
)

// Run a long process that needs to be manually killed. Killing a process is
// accomplished by sending a true boolean value to stopChannel. Logs all errors
// and output to the console.
//
// Internally, it sprouts three new goroutines. Two to read from os.Stdout and
// os.Stderr respectively, and a last one to kill the process when stopChannel
// receives true.
func Run(runCommand string, stopChannel chan bool) bool {
	cmd := exec.Command(runCommand)

	stdoutIn, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Error with StdoutPipe, error: %s\n", err.Error())
	}

	stderrIn, err := cmd.StderrPipe()
	if err != nil {
		log.Fatalf("Error with StderrPipe, error: %s\n", err.Error())
	}

	var stdoutLog, stderrLog logger.LogWriter
	go func() { io.Copy(stdoutLog, stdoutIn) }()
	go func() { io.Copy(stderrLog, stderrIn) }()

	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start command %s, error: %s\n", runCommand, err.Error())
	}

	go func() {
		stop := <-stopChannel
		if stop {
			pid := cmd.Process.Pid
			if err := cmd.Process.Kill(); err != nil {
				log.Fatalf("Killing process pidp %d failed, error: %s\n", pid, err.Error())
			}
		}
	}()

	return true
}
