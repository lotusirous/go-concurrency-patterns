package main

import "time"

type Item struct{ Tile, Channel, GUID string }

// Fetcher fetches Items and returns the time when the next fetch
// should be attempted. On failure, fetch returns non-nil error.
type Fetcher interface {
	Fetch() (item []Item, next time.Time, err error)
}

// A Subscription delivers Items over a channel.  Close cancels the
// subscription, closes the Updates channel, and returns the last fetch error,
// if any.
type Subscription interface {
	Updates() <-chan Item
	Close() error
}

func Subscribe(fetcher Fetcher) Subscription {
	s := &sub{
		fetcher: fetcher,
		updates: make(chan Item),
		closing: make(chan chan error),
	}
	return s
}

type sub struct {
	fetcher Fetcher
	updates chan Item
	closing chan chan error
}

func (s *sub) Updates() <-chan Item {
	return s.updates
}

func (s *sub) Close() error {
	// STOPCLOSESIG OMIT
	errc := make(chan error)
	s.closing <- errc // HLchan
	return <-errc     // HLchan
}

func main() {

}
