package watchman

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
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

func Connect(sockname string) (net.Conn, error) {
	conn, err := net.Dial("unix", sockname)
	return conn, err
}

// Subscribe to watchman's messages about a directory's filesystem changes,
// using an appropriate Unix Socket to watchman. Can take in the directory and
// a subscription name to configure the subscription.
//
// Deprecated!
func Subscribe(conn net.Conn, directory, subname string, suffixes []string) {
	var cmd string

	suffixesString := formatSuffixes(suffixes)
	cmd = fmt.Sprintf(`["subscribe", "%s", "%s", {"expression": ["allof", ["type", "f"], ["not", "empty"]%s ], "fields": ["name"] }]`,
		directory, subname, suffixesString)
	fmt.Println(cmd)

	_, err := conn.Write([]byte(cmd + "\n"))
	if err != nil {
		log.Fatalf("Error writing to socket, error: %s\n", err.Error())
	}
}

// WatchProject watches a project directory and all it's subdirectories.
// Currectly, this function is deprecated in favor of Subscribe, until
// further notice.
func WatchProject(conn net.Conn, directory string) {
	var cmd string

	cmd = fmt.Sprintf(`["watch", "%s"]`, directory)
	bytes, err := conn.Write([]byte(cmd + "\n"))
	if err != nil {
		log.Fatalf("Error writing to socket, error: %s\n", err.Error())
	}
	log.Printf("Bytes written: %d\n", bytes)
}

func formatSuffixes(suffixes []string) string {
	if len(suffixes) == 0 {
		return ""
	}

	suffixString := `, ["anyof"`

	for _, suffix := range suffixes {
		suffixString += fmt.Sprintf(`, ["suffix", "%s"]`, suffix)
	}
	return suffixString + `]`
}
