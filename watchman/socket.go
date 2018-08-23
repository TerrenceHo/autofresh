package watchman

import (
	"bytes"
	"encoding/json"
	"log"
	"os/exec"
)

// WatchmanSockName convert data from watchman get-sockname.
type WatchmanSockName struct {
	Version  string `json:"version"`
	Sockname string `json:"sockname"`
}

// GetSockName gets the socket filename using exec "watchman get-sockname".
// Passes stderr and stdout of command to the console. If it fails, then
// program shuts down. Otherwise, returns sockname.
func GetSockName(watchmanPath string) string {
	cmd := exec.Command("watchman", "get-sockname")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("Watchman failed to get sockname, error: %s\n", err.Error())
	}

	var wsn WatchmanSockName
	if err := json.NewDecoder(&stdout).Decode(&wsn); err != nil {
		log.Fatalf("Socket name did not decode correctly. Error: %s\n", err.Error())
	}

	return wsn.Sockname
}
