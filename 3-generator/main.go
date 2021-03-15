package main

import (
	"fmt"
	"math/rand"
	"time"
)

// boring is a function that returns a channel to communicate with it.
// <-chan string means receives-only channel of string.
func boring(msg string) <-chan string {
	c := make(chan string)
	// we launch goroutine inside a function
	// that sends the data to channel
	go func() {
		// The for loop simulate the infinite sender.
		for i := 0; i < 10; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}

		// The sender should close the channel
		close(c)

	}()
	return c // return a channel to caller.
}

func main() {

	joe := boring("Joe")
	ahn := boring("Ahn")

	// This loop yields 2 channels in sequence
	for i := 0; i < 10; i++ {
		fmt.Println(<-joe)
		fmt.Println(<-ahn)
	}

	// or we can simply use the for range
	// for msg := range joe {
	// 	fmt.Println(msg)
	// }
	fmt.Println("You're both boring. I'm leaving")

}
