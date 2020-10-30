package main

import (
	"fmt"
	"math/rand"
	"time"
)

// now, the boring function additional parametter
// `c chan string` is a channel
func boring(msg string, c chan string) {
	for i := 0; ; i++ {
		// send the value to channel
		// it also waits for receiver to be ready
		c <- fmt.Sprintf("%s %d", msg, i)
		time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
	}
}

func main() {

	// init our channel
	c := make(chan string)
	// placeholdering our goroutine
	go boring("boring!", c)

	for i := 0; i < 5; i++ {
		// <-c read the value from boring function
		// <-c waits for a value to be sent
		fmt.Printf("You say: %q\n", <-c)
	}
	fmt.Println("You're boring. I'm leaving")

}
