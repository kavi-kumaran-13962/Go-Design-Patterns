package main

import (
	"sync"
)

// logger is a simple struct for logging.
type logger struct {
}

var (
	loggerInstance *logger   // The single instance of the logger.
	once           sync.Once // A sync.Once ensures initialization is done only once.
)

// GetInstance returns the single instance of the logger.
func GetInstance() *logger {
	// Use sync.Once to guarantee that loggerInstance is initialized only once.
	once.Do(func() {
		println("Logger instance created")
		loggerInstance = &logger{}
	})
	return loggerInstance
}

// Log prints a log message.
func (lg logger) Log(msg string) {
	println(msg)
}

func main() {
	logger := GetInstance()
	logger.Log("Print msg 1")

	logger2 := GetInstance()
	logger2.Log("Print msg 2")
}
