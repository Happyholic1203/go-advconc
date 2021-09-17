package rss

import (
	"fmt"
	"math/rand"
	"time"
)

type Item struct {
	Title, Channel, GUID string
}

type Fetcher interface {
	Fetch() (items []Item, next time.Time, err error)
}

type Subscription interface {
	Updates() <-chan Item
	Close() error
}

func Fetch(domain string) Fetcher {
	return &fakeFetcher{channel: domain}
}

type fakeFetcher struct {
	channel  string
	numItems int
}

func (f *fakeFetcher) Fetch() (items []Item, next time.Time, err error) {
	next = time.Now().Add(time.Duration(rand.Intn(5)) * 500 * time.Millisecond)
	item := Item{
		Channel: f.channel,
		Title:   fmt.Sprintf("Item %d", f.numItems),
	}
	item.GUID = item.Channel + "/" + item.Title
	f.numItems++
	return []Item{item}, next, nil
}
