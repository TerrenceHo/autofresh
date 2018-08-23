package autofresh

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"path/filepath"

	"github.com/TerrenceHo/autofresh/config"
	"github.com/TerrenceHo/autofresh/runner"
	"github.com/TerrenceHo/autofresh/watchman"
)

const logo = `
    ___         __        ______               __  
   /   | __  __/ /_____  / ____/_______  _____/ /_ 
  / /| |/ / / / __/ __ \/ /_  / ___/ _ \/ ___/ __ \
 / ___ / /_/ / /_/ /_/ / __/ / /  /  __(__  ) / / /
/_/  |_\__,_/\__/\____/_/   /_/   \___/____/_/ /_/ 
`

const (
	subscriptionName = "test_autofresh"
)

var (
	startChannel    chan bool
	stopChannel     chan bool
	watchmanCommand string
	buildCommand    string
	runCommand      string
	suffixes        []string
)

// Start application. Will check if watchman exists at that path,
// retrieve the socket name, instantiate a connection to watchman using Unix
// Sockets, subscribe to the directory, and begin reading the subscription
// messages and building the executable.
func Start(conf config.Config) {
	if !conf.Hidebanner {
		fmt.Println(logo)
	}
	startChannel = make(chan bool)
	stopChannel = make(chan bool)
	buildCommand = conf.Build
	runCommand = conf.Run
	watchmanCommand = conf.Watchman
	suffixes = conf.Suffixes
	isRunning := false

	directory, err := filepath.Abs("./")
	if err != nil {
		log.Fatalf("Directory absolute file path cannot be found. Error: %s\n", err.Error())
	}
	watchman.Check(watchmanCommand)

	sockname := watchman.GetSockName(watchmanCommand)

	conn, err := net.Dial("unix", sockname)
	if err != nil {
		log.Fatalf("Dialing unix socket %s failed, error: %s\n",
			sockname, err.Error())
	}
	defer conn.Close()

	go read(conn, startChannel) // read in separate go routine

	// watchman.WatchProject(conn, directory)
	watchman.Subscribe(conn, directory, subscriptionName, suffixes)

	for {
		start := <-startChannel
		if start {
			// kill previous process if necessary
			if isRunning {
				stopChannel <- true
			}

			buildErr, err := runner.Build(buildCommand)
			if err != nil {
				fmt.Println(buildErr)
				fmt.Println(err.Error())
				continue // skip run if build returned with error
			}

			// run here
			isRunning = runner.Run(runCommand, stopChannel)
		}
	}
}

func read(c net.Conn, startChannel chan bool) {
	d := json.NewDecoder(c)
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
