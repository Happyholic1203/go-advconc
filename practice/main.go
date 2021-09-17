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
	// TODO
}

func Subscribe(fetcher rss.Fetcher) rss.Subscription {
	return &sub{}
}

func (s *sub) Updates() <-chan rss.Item {
	// TODO
	return nil
}

func (s *sub) Close() error {
	// TODO
	return nil
}

type mergedSub struct {
	// TODO
}

func Merge(subs ...rss.Subscription) rss.Subscription {
	// TODO
	return nil
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
