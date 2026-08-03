package main

import (
	"sync"
	"testing"
	"testing/synctest"
)

// TestRunDoesNotDeadlock guards Run's overflow path.
//
// The old code dropped the oldest item with a blocking receive (<-r.outCh).
// Finding outCh full and popping from it are two separate steps, so a consumer
// can empty outCh in between, leaving Run blocked forever on a value that only
// Run itself can send.
//
// synctest fails the moment every goroutine in the bubble is blocked, so this
// reports a deadlock instead of hanging. Capacity 1 keeps outCh full, so Run
// almost always takes the overflow path, and the rounds give the consumer
// enough chances to win the race.
func TestRunDoesNotDeadlock(t *testing.T) {
	const (
		rounds = 200
		items  = 2000
	)

	for round := 0; round < rounds; round++ {
		synctest.Test(t, func(t *testing.T) {
			inCh := make(chan int)
			outCh := make(chan int, 1)
			rb := NewRingBuffer(inCh, outCh)

			var wg sync.WaitGroup
			wg.Go(rb.Run)
			wg.Go(func() {
				for range outCh { // race Run to empty the buffer
				}
			})

			for i := 0; i < items; i++ {
				inCh <- i
			}
			close(inCh)
			wg.Wait()
		})
	}
}

// TestRunDropsOldest checks that a full buffer drops the oldest item rather
// than the newest. synctest.Wait returns once Run has finished and closed
// outCh, so the drain below is deterministic.
func TestRunDropsOldest(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		inCh := make(chan int)
		outCh := make(chan int, 4)
		rb := NewRingBuffer(inCh, outCh)
		go rb.Run()

		for i := 0; i < 10; i++ {
			inCh <- i
		}
		close(inCh)
		synctest.Wait() // Run has finished and closed outCh

		var got []int
		for v := range outCh {
			got = append(got, v)
		}

		want := []int{6, 7, 8, 9}
		if len(got) != len(want) {
			t.Fatalf("got %v (%d items), want %v (%d items)", got, len(got), want, len(want))
		}
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("got %v, want %v", got, want)
			}
		}
	})
}
