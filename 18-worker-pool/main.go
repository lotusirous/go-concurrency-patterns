// Credit:
// https://gobyexample.com/worker-pools

package main

import (
	"fmt"
	"time"
)

func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("worker", id, "started job", j)
		time.Sleep(time.Second)
		fmt.Println("worker", id, "fnished job", j)
		results <- j * 2
	}
}

func main() {
	const numbJobs = 8
	jobs := make(chan int, numbJobs)
	results := make(chan int, numbJobs)

	// start a group of workers
	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}

	// send data to worker group
	for j := 1; j <= numbJobs; j++ {
		jobs <- j
	}
	close(jobs)
	fmt.Println("Closed job")
	for a := 1; a <= numbJobs; a++ {
		<-results
	}

}
