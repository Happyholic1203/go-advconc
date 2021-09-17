# Advanced Go Concurrency Patterns

A simple (mock) RSS subscription fetcher for practicing go concurrency patterns:
- Use for-select loop to write more responsive code
- Use `nil` channels to disable select cases
- Use `chan chan error` to avoid data race that usually happen in `Close()`

The directory `practice` contains the necessary boilerplate code for you to practice on your own.

To give it a run: `cd improved; go run main.go`

To practice on your own:
1. `cd practice`
2. `vim main.go`
3. `go run main.go`: go to 2 if anything goes wrong

All credit goes to Sameer Ajmani and his great talk *Advanced Go Concurrency Patterns*:
- [Slides](https://talks.golang.org/2013/advconc.slide#1)
- [Youtube](https://www.youtube.com/watch?v=QDDwwePbDtw)