package main

import (
	"sync"
)

// producer function takes a variadic parameter of integers and returns a receive-only directional channel of integers
func producer(nums ...int) <-chan int {
	myChannel := make(chan int) // declare a channel

	go func() {
		// iterate the nums data and sends it to channel
		for _, val := range nums {
			myChannel <- val
		}
		close(myChannel)
	}()

	return myChannel
}

func split(ch <-chan int, n int) []chan int {
	var channels []chan int
	var wg sync.WaitGroup
	wg.Add(1)

	for i := 0; i < n; i++ {
		channels = append(channels, make(chan int))
	}

	go func() {
		defer wg.Done()
		for no := range ch {
			for _, c := range channels {
				c <- no
			}
		}

	}()

	go func() {
		wg.Wait()
		for _, c := range channels {
			close(c)
		}
	}()

	return channels
}

func main() {
	data1 := []int{1, 2, 3, 4, 5}

	//it receives a "receive-only" directional channel
	ch1 := producer(data1...)
	var wg sync.WaitGroup

	out := split(ch1, 3)
	for _, ch := range out {
		wg.Add(1)
		go func(c <-chan int) {
			defer wg.Done()
			for no := range c {
				println(no) // print the integers received from the output channel
			}
		}(ch)
	}

	wg.Wait() // wait for all goroutines to complete
}
