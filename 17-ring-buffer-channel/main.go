package main

import "log"

// A channel-based ring buffer removes the oldest item when the queue is full
// Ref:
// https://tanzu.vmware.com/content/blog/a-channel-based-ring-buffer-in-go

func NewRingBuffer(inCh, outCh chan int) *ringBuffer {
	return &ringBuffer{
		inCh:  inCh,
		outCh: outCh,
	}
}

// ringBuffer throttle buffer for implement async channel.
type ringBuffer struct {
	inCh  chan int
	outCh chan int
}

// Run moves items from inCh to outCh, dropping the oldest item when outCh is
// full, and closes outCh once inCh is closed.
//
// Run must be the only sender on outCh.
func (r *ringBuffer) Run() {
	for v := range r.inCh {
		// There is room: just send.
		select {
		case r.outCh <- v:
			continue
		default:
		}

		// outCh looked full, so drop the oldest item. Never block here: a
		// consumer may have emptied outCh since the select above, and a
		// blocking receive would wait for a value only Run itself can send.
		select {
		case <-r.outCh:
		default:
		}

		// Safe: Run is the only sender, so the slot freed above is still free.
		r.outCh <- v
	}
	close(r.outCh)
}

func main() {
	inCh := make(chan int)
	outCh := make(chan int, 4) // try to change outCh buffer to understand the result
	rb := NewRingBuffer(inCh, outCh)
	go rb.Run()

	for i := 0; i < 10; i++ {
		inCh <- i
	}

	close(inCh)

	for res := range outCh {
		log.Println(res)
	}

}
