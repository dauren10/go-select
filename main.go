package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// A channel for communication between the consumer and processor workers
var data = make(chan int)
var processed = make(chan int)
var unblock = make(chan bool)
var unblock_flag bool

// A wait group to ensure all workers have finished before exiting the program
var wg sync.WaitGroup

// The consumer worker reads data from the data channel and sends it to the processor
func consumer() {
	defer wg.Done()

	// Process data until the data channel is closed
	for {
		select {
		case d, ok := <-data:
			if !ok {
				// Data channel is closed, exit the loop
				return
			}
			fmt.Println("Received data:", d)
			processed <- d
		case <-time.After(time.Second):
			// Pause for 1 second if no data is available on the data channel
			fmt.Println("Pausing...")
		case <-unblock:
			// Unblock the loop if an unblock signal is received

			fmt.Println("Unblocking...")
		}
	}
}

// The processor worker reads data from the processed channel and processes it
func processor() {
	defer wg.Done()
	unblock_flag = true
	// Process data until the processed channel is closed
	for i := 0; i < 30; i++ {
		d := rand.Intn(1000)
		fmt.Println("Sending data:", d)
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		println(unblock_flag)

	}
}

func main() {
	// Add two workers to the wait group
	wg.Add(2)

	// Start the consumer and processor workers
	go consumer()
	go processor()

	// Generate and send 5 random integers to the data channel
	for i := 0; i < 30; i++ {
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		if i > 3 {
			unblock <- true
		}
	}
	// Unblock the consumer worker by sending a signal on the unblock channel

	// Close the data channel to signal the consumer worker to stop iterating
	close(data)

	// Wait for all workers to finish
	wg.Wait()
}
