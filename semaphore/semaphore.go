package main

import (
	"fmt"
	"time"
)

// semaphore interface with Acquire and Release methods
type semaphore interface {
	Acquire()
	Release()
}

// semaphoreImpl struct with a channel of struct{} type
type semaphoreImpl struct {
	sem chan struct{}
}

// NewSemaphore returns a new semaphore with a channel of struct{} type
func NewSemaphore(maxConcurrent int) semaphore {
	return &semaphoreImpl{
		sem: make(chan struct{}, maxConcurrent),
	}
}

// Acquire acquires a semaphore
func (s *semaphoreImpl) Acquire() {
	s.sem <- struct{}{}
}

// Release releases a semaphore
func (s *semaphoreImpl) Release() {
	<-s.sem
}

// longRunningProcess simulates a long running process
func longRunningProcess(taskID int) {
	fmt.Println(
		time.Now().Format("15:04:05"),
		"Running task with ID",
		taskID)
	time.Sleep(2 * time.Second)
}

func main() {
	totalTasks := 10
	totalWorkers := 2
	sem := NewSemaphore(totalWorkers)
	doneC := make(chan bool, 1)

	// loop through the total number of tasks
	for i := 0; i < totalTasks; i++ {
		sem.Acquire() // acquire a semaphore
		go func(i int) {
			defer sem.Release() // release a semaphore
			longRunningProcess(i)
			if i == totalTasks-1 {
				doneC <- true // send a signal to the done channel
			}
		}(i)
	}

	<-doneC // wait for the done signal
}
