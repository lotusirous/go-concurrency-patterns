package main

import (
	"fmt"
	"math/rand"
	"time"
)

// the boring function return a channel to communicate with it.
func boring(msg string) <-chan string { // <-chan string means receives-only channel of string.
	c := make(chan string)
	go func() { // we launch goroutine inside a function.
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}

	}()
	return c // return a channel to caller.
}

func main() {

	// Now, we don't need a this channel anymore
	// c := make(chan string)
	joe := boring("Joe") // function returns a channel
	ahn := boring("Ahn")

	for i := 0; i < 5; i++ {
		fmt.Println(<-joe)
		fmt.Println(<-ahn)
	}
	fmt.Println("You're both boring. I'm leaving")

}
