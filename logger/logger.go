package logger

import (
	"fmt"
	"log"
)

// LogWriter contains a logger, to be used with logging to the console.
type LogWriter struct {
	logger *log.Logger
}

// Write implements io.Writer interface. Allows loggers to write through to the
// console.
func (lw LogWriter) Write(p []byte) (n int, err error) {
	fmt.Print(string(p))
	return len(p), nil
}
