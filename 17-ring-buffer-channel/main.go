package main

import "log"

// A channel-based ring buffer removes the oldest item when the queue is full
// Ref:
// https://tanzu.vmware.com/content/blog/a-channel-based-ring-buffer-in-go

func NewRingBuffer(inCh, outCh chan int, down chan bool) *ringBuffer {
	return &ringBuffer{
		inCh:  inCh,
		outCh: outCh,
		down:  down,
	}
}

// ringBuffer throttle buffer for implement async channel.
type ringBuffer struct {
	inCh  chan int
	outCh chan int
	down  chan bool
}

func (r *ringBuffer) Run() {
	for v := range r.inCh {
		select {
		case r.outCh <- v:
		default:
			<-r.outCh // pop one item from outchan
			r.outCh <- v
		}
	}
	r.down <- true
	close(r.outCh)
}

func main() {
	inCh := make(chan int)
	outCh := make(chan int, 4) // try to change outCh buffer to understand the result
	down := make(chan bool, 1)
	rb := NewRingBuffer(inCh, outCh, down)
	go rb.Run()

	for i := 0; i < 10; i++ {
		inCh <- i
	}

	close(inCh)

	if <-down == true {
		for res := range outCh {
			log.Println(res)
		}
	}
}
