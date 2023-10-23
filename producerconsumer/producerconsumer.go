package main

import (
	"fmt"
	"sync"
	"time"
)

// Producer is responsible for producing data and sending it to the channel.
func Producer(ch chan int) {
	for i := 0; i < 20; i++ {
		ch <- i // Send data to the channel
		time.Sleep(100 * time.Millisecond)
	}
	close(ch) // Close the channel to signal no more data will be sent
}

// Consumer is responsible for consuming data from the channel.
func Consumer(ch chan int) {
	for data := range ch {
		fmt.Println("Consumed", data) // Print the consumed data
	}
}

func main() {
	// Create a buffered channel to hold data
	ch := make(chan int, 2)

	// Create a WaitGroup to wait for goroutines to complete
	var wg sync.WaitGroup

	// Add 2 to the WaitGroup to wait for two goroutines
	wg.Add(2)

	// Goroutine for the producer
	go func() {
		defer wg.Done() // Decrement the WaitGroup when the goroutine is done
		Producer(ch)
	}()

	// Goroutine for the consumer
	go func() {
		defer wg.Done() // Decrement the WaitGroup when the goroutine is done
		Consumer(ch)
	}()

	// Wait for both producer and consumer to finish
	wg.Wait()
}
