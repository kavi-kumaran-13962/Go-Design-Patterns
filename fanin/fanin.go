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

// fanIn function takes a variadic parameter of receive-only directional channels of integers and returns a receive-only directional channel of integers
func fanIn(inputs ...<-chan int) <-chan int {
	wg := sync.WaitGroup{} // declare a WaitGroup
	wg.Add(len(inputs))    // add the number of inputs to the WaitGroup
	output := make(chan int)

	// iterate over the inputs and start a goroutine for each input
	for _, in := range inputs {
		go func(ch <-chan int) {
			defer wg.Done() // mark the WaitGroup as done when the goroutine completes
			for n := range ch {
				output <- n // send the integer to the output channel
			}

		}(in)
	}

	go func() {
		wg.Wait()     // wait for all goroutines to complete
		close(output) // close the output channel
	}()
	return output
}

func main() {
	data1 := []int{1, 2, 3, 4, 5}
	data2 := []int{10, 20, 30, 40, 50}
	var wg sync.WaitGroup

	//it receives a "receive-only" directional channel
	ch1 := producer(data1...)
	ch2 := producer(data2...)

	out := fanIn(ch1, ch2) // combine the two channels using fanIn function
	for no := range out {
		println(no) // print the integers received from the output channel
	}
	wg.Wait() // wait for all goroutines to complete
}
