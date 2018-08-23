package watchman

import (
	"log"
	"os/exec"
)

// Check that watchman exists in the system. If not error out and shut down.
// Otherwise, continue execution.
func Check(watchmanPath string) {
	_, err := exec.LookPath(watchmanPath)
	if err != nil {
		log.Printf("Failed to find watchman executable. Error: %s\n", err.Error())
		log.Fatal("Please install watchman before using autofresh.")
	}
}
