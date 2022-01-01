package main

import (
	"fmt"
	"math/rand"
	"time"
)

// An Item is a stripped-down RSS item.
type Item struct{ Title, Channel, GUID string }

// A Fetcher fetches Items and returns the time when the next fetch should be
// attempted.  On failure, Fetch returns a non-nil error.
type Fetcher interface {
	Fetch() (items []Item, next time.Time, err error)
}

// A Subscription delivers Items over a channel.  Close cancels the
// subscription, closes the Updates channel, and returns the last fetch error,
// if any.
type Subscription interface {
	Updates() <-chan Item
	Close() error
}

// Subscribe returns a new Subscription that uses fetcher to fetch Items.
func Subscribe(fetcher Fetcher) Subscription {
	s := &sub{
		fetcher: fetcher,
		updates: make(chan Item),       // for Updates
		closing: make(chan chan error), // for Close
	}
	go s.loop()
	return s
}

// sub implements the Subscription interface.
type sub struct {
	fetcher Fetcher         // fetches items
	updates chan Item       // sends items to the user
	closing chan chan error // for Close
}

func (s *sub) Updates() <-chan Item {
	return s.updates
}

func (s *sub) Close() error {
	errc := make(chan error)
	s.closing <- errc // HLchan
	return <-errc     // HLchan
}

// loopCloseOnly is a version of loop that includes only the logic
// that handles Close.
func (s *sub) loopCloseOnly() {
	var err error // set when Fetch fails
	for {
		select {
		case errc := <-s.closing: // HLchan
			errc <- err      // HLchan
			close(s.updates) // tells receiver we're done
			return
		}
	}

}

// loopFetchOnly is a version of loop that includes only the logic
// that calls Fetch.
func (s *sub) loopFetchOnly() {
	var (
		pending []Item    // appended by fetch; consumed by send
		next    time.Time // initially January 1, year 0
		err     error
	)
	for {
		var fetchDelay time.Duration // initally 0 (no delay)
		if now := time.Now(); next.After(now) {
			fetchDelay = next.Sub(now)
		}
		startFetch := time.After(fetchDelay)

		select {
		case <-startFetch:
			var fetched []Item
			fetched, next, err = s.fetcher.Fetch()
			if err != nil {
				next = time.Now().Add(10 * time.Second)
				break
			}
			pending = append(pending, fetched...)
		}
	}

}

// loopSendOnly is a version of loop that includes only the logic for
// sending items to s.updates.
func (s *sub) loopSendOnly() {

	var pending []Item // appended by fetch; consumed by send
	for {
		var first Item
		var updates chan Item // HLupdates
		if len(pending) > 0 {
			first = pending[0]
			updates = s.updates // enable send case // HLupdates
		}

		select {
		case updates <- first:
			pending = pending[1:]
		}
	}

}

// mergedLoop is a version of loop that combines loopCloseOnly,
// loopFetchOnly, and loopSendOnly.
func (s *sub) mergedLoop() {

	var (
		pending []Item
		next    time.Time
		err     error
	)

	for {
		var fetchDelay time.Duration
		if now := time.Now(); next.After(now) {
			fetchDelay = next.Sub(now)
		}
		startFetch := time.After(fetchDelay)

		var first Item
		var updates chan Item
		if len(pending) > 0 {
			first = pending[0]
			updates = s.updates // enable send case
		}

		select {
		case errc := <-s.closing: // HLcases
			errc <- err
			close(s.updates)
			return

		case <-startFetch: // HLcases
			var fetched []Item
			fetched, next, err = s.fetcher.Fetch() // HLfetch
			if err != nil {
				next = time.Now().Add(10 * time.Second)
				break
			}
			pending = append(pending, fetched...) // HLfetch

		case updates <- first: // HLcases
			pending = pending[1:]
		}

	}
}

// dedupeLoop extends mergedLoop with deduping of fetched items.
func (s *sub) dedupeLoop() {
	const maxPending = 10

	var pending []Item
	var next time.Time
	var err error
	var seen = make(map[string]bool) // set of item.GUIDs // HLseen

	for {

		var fetchDelay time.Duration
		if now := time.Now(); next.After(now) {
			fetchDelay = next.Sub(now)
		}
		var startFetch <-chan time.Time // HLcap
		if len(pending) < maxPending {  // HLcap
			startFetch = time.After(fetchDelay) // enable fetch case  // HLcap
		} // HLcap

		var first Item
		var updates chan Item
		if len(pending) > 0 {
			first = pending[0]
			updates = s.updates // enable send case
		}
		select {
		case errc := <-s.closing:
			errc <- err
			close(s.updates)
			return

		case <-startFetch:
			var fetched []Item
			fetched, next, err = s.fetcher.Fetch() // HLfetch
			if err != nil {
				next = time.Now().Add(10 * time.Second)
				break
			}
			for _, item := range fetched {
				if !seen[item.GUID] { // HLdupe
					pending = append(pending, item) // HLdupe
					seen[item.GUID] = true          // HLdupe
				} // HLdupe
			}

		case updates <- first:
			pending = pending[1:]
		}
	}
}

// loop periodically fecthes Items, sends them on s.updates, and exits
// when Close is called.  It extends dedupeLoop with logic to run
// Fetch asynchronously.
func (s *sub) loop() {
	const maxPending = 10
	type fetchResult struct {
		fetched []Item
		next    time.Time
		err     error
	}

	var fetchDone chan fetchResult // if non-nil, Fetch is running // HL

	var pending []Item
	var next time.Time
	var err error
	var seen = make(map[string]bool)
	for {
		var fetchDelay time.Duration
		if now := time.Now(); next.After(now) {
			fetchDelay = next.Sub(now)
		}

		var startFetch <-chan time.Time
		if fetchDone == nil && len(pending) < maxPending { // HLfetch
			startFetch = time.After(fetchDelay) // enable fetch case
		}

		var first Item
		var updates chan Item
		if len(pending) > 0 {
			first = pending[0]
			updates = s.updates // enable send case
		}

		select {
		case <-startFetch: // HLfetch
			fetchDone = make(chan fetchResult, 1) // HLfetch
			go func() {
				fetched, next, err := s.fetcher.Fetch()
				fetchDone <- fetchResult{fetched, next, err}
			}()
		case result := <-fetchDone: // HLfetch
			fetchDone = nil // HLfetch
			// Use result.fetched, result.next, result.err

			fetched := result.fetched
			next, err = result.next, result.err
			if err != nil {
				next = time.Now().Add(10 * time.Second)
				break
			}
			for _, item := range fetched {
				if id := item.GUID; !seen[id] { // HLdupe
					pending = append(pending, item)
					seen[id] = true // HLdupe
				}
			}
		case errc := <-s.closing:
			errc <- err
			close(s.updates)
			return
		case updates <- first:
			pending = pending[1:]
		}
	}
}

// naiveMerge is a version of Merge that doesn't quite work right.  In
// particular, the goroutines it starts may block forever on m.updates
// if the receiver stops receiving.
type naiveMerge struct {
	subs    []Subscription
	updates chan Item
}

func NaiveMerge(subs ...Subscription) Subscription {
	m := &naiveMerge{
		subs:    subs,
		updates: make(chan Item),
	}

	for _, sub := range subs {
		go func(s Subscription) {
			for it := range s.Updates() {
				m.updates <- it // HL
			}
		}(sub)
	}

	return m
}

func (m *naiveMerge) Close() (err error) {
	for _, sub := range m.subs {
		if e := sub.Close(); err == nil && e != nil {
			err = e
		}
	}
	close(m.updates) // HL
	return
}

func (m *naiveMerge) Updates() <-chan Item {
	return m.updates
}

type merge struct {
	subs    []Subscription
	updates chan Item
	quit    chan struct{}
	errs    chan error
}

// Merge returns a Subscription that merges the item streams from subs.
// Closing the merged subscription closes subs.
func Merge(subs ...Subscription) Subscription {

	m := &merge{
		subs:    subs,
		updates: make(chan Item),
		quit:    make(chan struct{}),
		errs:    make(chan error),
	}

	for _, sub := range subs {
		go func(s Subscription) {
			for {
				var it Item
				select {
				case it = <-s.Updates():
				case <-m.quit: // HL
					m.errs <- s.Close() // HL
					return              // HL
				}
				select {
				case m.updates <- it:
				case <-m.quit: // HL
					m.errs <- s.Close() // HL
					return              // HL
				}
			}
		}(sub)
	}

	return m
}

func (m *merge) Updates() <-chan Item {
	return m.updates
}

func (m *merge) Close() (err error) {
	close(m.quit) // HL
	for range m.subs {
		if e := <-m.errs; e != nil { // HL
			err = e
		}
	}
	close(m.updates) // HL
	return
}

// NaiveDedupe converts a stream of Items that may contain duplicates
// into one that doesn't.
func NaiveDedupe(in <-chan Item) <-chan Item {
	out := make(chan Item)
	go func() {
		seen := make(map[string]bool)
		for it := range in {
			if !seen[it.GUID] {
				// BUG: this send blocks if the
				// receiver closes the Subscription
				// and stops receiving.
				out <- it // HL
				seen[it.GUID] = true
			}
		}
		close(out)
	}()
	return out
}

type deduper struct {
	s       Subscription
	updates chan Item
	closing chan chan error
}

// Dedupe converts a Subscription that may send duplicate Items into
// one that doesn't.
func Dedupe(s Subscription) Subscription {
	d := &deduper{
		s:       s,
		updates: make(chan Item),
		closing: make(chan chan error),
	}
	go d.loop()
	return d
}

func (d *deduper) loop() {
	in := d.s.Updates() // enable receive
	var pending Item
	var out chan Item // disable send
	seen := make(map[string]bool)
	for {
		select {
		case it := <-in:
			if !seen[it.GUID] {
				pending = it
				in = nil        // disable receive
				out = d.updates // enable send
				seen[it.GUID] = true
			}
		case out <- pending:
			in = d.s.Updates() // enable receive
			out = nil          // disable send
		case errc := <-d.closing:
			err := d.s.Close()
			errc <- err
			close(d.updates)
			return
		}
	}
}

func (d *deduper) Close() error {
	errc := make(chan error)
	d.closing <- errc
	return <-errc
}

func (d *deduper) Updates() <-chan Item {
	return d.updates
}

// Fetch returns a Fetcher for Items from domain.
func Fetch(domain string) Fetcher {
	return fakeFetch(domain)
}

func fakeFetch(domain string) Fetcher {
	return &fakeFetcher{channel: domain}
}

type fakeFetcher struct {
	channel string
	items   []Item
}

// FakeDuplicates causes the fake fetcher to return duplicate items.
var FakeDuplicates bool

func (f *fakeFetcher) Fetch() (items []Item, next time.Time, err error) {
	now := time.Now()
	next = now.Add(time.Duration(rand.Intn(5)) * 500 * time.Millisecond)
	item := Item{
		Channel: f.channel,
		Title:   fmt.Sprintf("Item %d", len(f.items)),
	}
	item.GUID = item.Channel + "/" + item.Title
	f.items = append(f.items, item)
	if FakeDuplicates {
		items = f.items
	} else {
		items = []Item{item}
	}
	return
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {

	// Subscribe to some feeds, and create a merged update stream.
	merged := Merge(
		Subscribe(Fetch("blog.golang.org")),
		Subscribe(Fetch("googleblog.blogspot.com")),
		Subscribe(Fetch("googledevelopers.blogspot.com")),
	)

	// Close the subscriptions after some time.
	time.AfterFunc(3*time.Second, func() {
		fmt.Println("closed:", merged.Close())
	})

	// Print the stream.
	for it := range merged.Updates() {
		fmt.Println(it.Channel, it.Title)
	}

	panic("show me the stacks")
}
