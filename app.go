package autofresh

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/TerrenceHo/autofresh/config"
	"github.com/TerrenceHo/autofresh/watcher"
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
	// isRunning := false

	// watchman.NewWatchman(watchmanCommand,

	directory, err := filepath.Abs("./")
	if err != nil {
		log.Fatalf("Directory absolute file path cannot be found. Error: %s\n", err.Error())
	}

	wm := watchman.NewWatchman(watchmanCommand, directory, subscriptionName)
	defer wm.Close()

	builder := watcher.NewBuilder(wm, buildCommand, runCommand)

	wm.Subscribe(directory, suffixes)
	builder.Refresh()

	// go read(conn, startChannel) // read in separate goroutine

}

// func read(c net.Conn, startChannel chan bool) {
// 	d := json.NewDecoder(c)
// 	m := make(map[string]interface{})
// 	var start = true
// 	for {
// 		if err := d.Decode(&m); err != nil {
// 			if err != io.EOF {
// 				log.Fatalf("Error decoding, error: %s\n", err.Error())
// 			}
// 		}
// 		fmt.Println(m)
// 		startChannel <- start
// 		start = !start
// 	}
// }
