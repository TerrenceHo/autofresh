package watcher

import (
	"fmt"

	"github.com/TerrenceHo/autofresh/runner"
)

// Watcher is an interface that defines how a file system watcher should
// operate. In a separate goroutine, it continuously reads the file system and
// marks everytime a file changes. It uses a boolean channel to indicate to
// another function when a file has changed, and what actions to take then.
type Watcher interface {
	Read(startChannel chan bool)
}

// A builder takes in a watcher, and when the watcher indicates a file has
// changed, it builds and runs a command.
type Builder struct {
	// watcher interface object.
	watcher Watcher

	// startChannel is used by the watcher to indicate when a file has changed.
	startChannel chan bool

	// communicate to a running process when a running program should be
	// stopped.
	stopChannel chan bool

	// A command to execute a command that terminates
	buildCommand string

	// A command to run a command that will not terminate, but must be
	// terminated by stopChannel
	runCommand string
}

// NewBuilder takes the appropriate commands and returns a Builder pointer
func NewBuilder(watcher Watcher, buildCommand, runCommand string) *Builder {
	return &Builder{
		watcher:      watcher,
		startChannel: make(chan bool),
		stopChannel:  make(chan bool),
		buildCommand: buildCommand,
		runCommand:   runCommand,
	}
}

func (b *Builder) Refresh() {
	go b.watcher.Read(b.startChannel)

	var isRunning bool = false
	for {
		start := <-b.startChannel
		if start {
			// kill previous process if necessary
			if isRunning {
				b.stopChannel <- true
			}

			if b.buildCommand != "" {
				buildErr, err := runner.Build(b.buildCommand)
				if err != nil {
					fmt.Println(buildErr)
					fmt.Println(err.Error())
					continue // skip run if build returned with error
				}
			}

			// run here
			if b.runCommand != "" {
				isRunning = runner.Run(b.runCommand, b.stopChannel)
			}
		}
	}
}
