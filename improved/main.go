package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/Happyholic1203/go-advconc/rss"
)

type sub struct {
	fetcher rss.Fetcher
	items   chan rss.Item
	closing chan chan error
}

func Subscribe(fetcher rss.Fetcher) rss.Subscription {
	s := &sub{fetcher, make(chan rss.Item), make(chan chan error)}
	go s.loop()
	return s
}

type fetchResult struct {
	items []rss.Item
	next  time.Time
	err   error
}

func (s *sub) loop() {
	var pendingItems []rss.Item
	var err error
	var next time.Time
	var pendingItem rss.Item
	seen := map[string]struct{}{}
	var fetchDone chan fetchResult
	for {
		var updates chan rss.Item
		if len(pendingItems) > 0 { // only enable updates when there are pending items
			updates = s.items
			pendingItem = pendingItems[0]
		}

		var startFetch <-chan time.Time
		const maxItems = 5 // to prevent allocating an arbitrary amount of memory
		if fetchDone == nil && len(pendingItems) < maxItems {
			// only enable fetch when we're not fetching and we still have room
			startFetch = time.After(next.Sub(time.Now()))
		}

		select {
		case errc := <-s.closing:
			errc <- err
			close(s.items)
			return
		case <-startFetch:
			fetchDone = make(chan fetchResult)
			go func() {
				items, next, err := s.fetcher.Fetch()
				fetchDone <- fetchResult{items, next, err}
			}()
		case result := <-fetchDone:
			fetchDone = nil
			next, err = result.next, result.err
			if err != nil {
				next = time.Now().Add(10 * time.Second)
				break
			}
			for _, item := range result.items {
				if _, ok := seen[item.Title]; !ok {
					seen[item.Title] = struct{}{}
					pendingItems = append(pendingItems, item)
				}
			}
		case updates <- pendingItem:
			pendingItems = pendingItems[1:]
		}
	}
}

func (s *sub) Updates() <-chan rss.Item {
	return s.items
}

func (s *sub) Close() error {
	errc := make(chan error)
	s.closing <- errc
	return <-errc
}

type mergedSub struct {
	subs    []rss.Subscription
	updates chan rss.Item
	stop    chan struct{}
}

func Merge(subs ...rss.Subscription) rss.Subscription {
	updates := make(chan rss.Item)
	stop := make(chan struct{})
	for _, s := range subs {
		go func(s rss.Subscription) {
			var update rss.Item
			for {
				select {
				case update = <-s.Updates():
				case <-stop:
					return
				}

				select {
				case updates <- update:
				case <-stop:
					return
				}
			}
		}(s)
	}
	return &mergedSub{subs, updates, stop}
}

func (s *mergedSub) Updates() <-chan rss.Item {
	return s.updates
}

func (s *mergedSub) Close() (err error) {
	close(s.stop)
	for _, s := range s.subs {
		if e := s.Close(); e != nil {
			err = e
		}
	}
	close(s.updates)
	return err
}

func main() {
	merged := Merge(
		Subscribe(rss.Fetch("blog.golang.org")),
		Subscribe(rss.Fetch("googleblog.blogspot.com")))

	var once sync.Once

	closeFunc := func() {
		fmt.Printf("Merged subscription closed. Errors: %v\n", merged.Close())
	}
	time.AfterFunc(3*time.Second, func() {
		once.Do(closeFunc)
	})

	userInterrupt := make(chan os.Signal)
	go func() {
		<-userInterrupt
		once.Do(closeFunc)
	}()
	signal.Notify(userInterrupt, os.Interrupt)

	for it := range merged.Updates() {
		fmt.Println(it.Channel, it.Title)
	}

	time.Sleep(time.Second)
}
