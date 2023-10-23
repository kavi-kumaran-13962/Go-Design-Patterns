package main

import (
	"sync"
	"time"
)

// Task interface defines the Execute method for all tasks.
type task interface {
	Execute()
}

// logTask represents a task that prints a message.
type logTask struct {
	msg string
}

// Execute method for logTask prints the message.
func (lg logTask) Execute() {
	println(lg.msg)
}

// doubleTask represents a task that prints a doubled number.
type doubleTask struct {
	no int
}

// Execute method for doubleTask prints the doubled number.
func (db doubleTask) Execute() {
	println(db.no * 2)
}

// worker is responsible for executing tasks from the task channel.
func worker(taskch chan task) {
	for task := range taskch {
		task.Execute()
	}
}

// taskProducer produces a batch of log and double tasks.
func taskProducer(taskChannel chan task) {
	for i := 0; i < 1000; i++ {
		logTask := logTask{"log"}
		taskChannel <- logTask
	}
	for i := 0; i < 1000; i++ {
		doubleTask := doubleTask{i}
		taskChannel <- doubleTask
	}
	close(taskChannel)
}

func main() {
	taskChannel := make(chan task, 2)

	// Try increasing the number of workers
	numWorkers := 20

	var wg sync.WaitGroup
	wg.Add(1)

	startTime := time.Now()

	// Launch the taskProducer goroutine.
	go func() {
		defer wg.Done()
		taskProducer(taskChannel)
	}()

	// Launch worker goroutines.
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(taskChannel)
		}()
	}
	wg.Wait()
	endTime := time.Now()

	duration := endTime.Sub(startTime)
	println("Execution time (milliseconds):", duration.Milliseconds())
}
