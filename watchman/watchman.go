package watchman

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
)

// Watchman is a struct that implements the Watcher Interface.
type Watchman struct {
	conn            net.Conn
	watchmanCommand string
	directory       string
	subname         string
}

// NewWatchman creates a new watchman object. NewWatchman will error out if the
// watchman command does not exist on the system, or if connecting to the
// watchman socket does not work.
func NewWatchman(watchmanCommand, directory, subname string) *Watchman {
	Check(watchmanCommand)
	sockname := GetSockName(watchmanCommand)
	conn, err := Connect(sockname)
	if err != nil {
		log.Fatalf("Dialing unix socket %s failed, error: %s\n",
			sockname, err.Error())
	}

	return &Watchman{
		conn:            conn,
		watchmanCommand: watchmanCommand,
		directory:       directory,
		subname:         subname,
	}
}

// Read takes in boolean channel and sends true to the channel everytime a new
// message is decoded.
func (wm *Watchman) Read(startChannel chan bool) {
	d := json.NewDecoder(wm.conn)
	m := make(map[string]interface{})
	var start = true
	for {
		if err := d.Decode(&m); err != nil {
			if err != io.EOF {
				log.Fatalf("Error decoding, error: %s\n", err.Error())
			}
		}
		fmt.Println(m)
		startChannel <- start
		start = !start
	}
}

func (wm *Watchman) Subscribe(directory string, suffixes []string) {
	var cmd string

	suffixesString := formatSuffixes(suffixes)
	cmd = fmt.Sprintf(`["subscribe", "%s", "%s", {"expression": ["allof", ["type", "f"], ["not", "empty"]%s ], "fields": ["name"] }]`,
		directory, wm.subname, suffixesString)
	fmt.Println(cmd)

	_, err := wm.conn.Write([]byte(cmd + "\n"))
	if err != nil {
		log.Fatalf("Error writing to socket, error: %s\n", err.Error())
	}

}

// Close implements the interface to close the connection. Allows program to
// defer closing the socket.
func (wm *Watchman) Close() {
	wm.conn.Close()
}
