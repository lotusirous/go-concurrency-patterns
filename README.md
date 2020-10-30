# Go Concurrency Patterns

This repository aims to collect pretty patterns for concurrency in Golang



## Basic

- Slide: https://talks.golang.org/2012/concurrency.slide#1

- Code: https://talks.golang.org/2012/concurrency/support/

## Advance

- Slide: https://talks.golang.org/2013/advconc.slide

## (2018) Rethinking classical concurrency pattern

- https://www.youtube.com/watch?v=5zXAHh5tJqQ

Principles:

1. Start goroutine when you have concurrent work.
2. Share by communication.


- Asynchronous: **returns** to the caller **before its result is ready**.
```go
func Fetch(name string) <-chan Item {
    c := make(chan Item, 1)
}
```
- Condition variables
- Worker pools
