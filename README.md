# Go Concurrency Patterns

This repository collects common concurrency patterns in Golang


## Materials
- [Concurrency is not parallelism](https://blog.golang.org/waza-talk)
- [Go Concurrency Patterns](https://talks.golang.org/2012/concurrency.slide#1)
and [source](https://talks.golang.org/2012/concurrency/support/)
- [Advanced Go Concurrency Patterns](https://talks.golang.org/2013/advconc.slide)
- [Rethinking classical concurrency pattern](https://www.youtube.com/watch?v=5zXAHh5tJqQ)
- [Go Concurrency Patterns: Pipelines and cancellation](https://blog.golang.org/pipelines)
- [Go Concurrency Patterns: Timing out, moving on](https://blog.golang.org/concurrency-timeouts)
- [Complex Concurrency Patterns with Go](https://www.youtube.com/watch?v=2HOO5gIgyMg)

## Context:
- [Go Concurrency Patterns: Context](https://blog.golang.org/context)
- [How to correctly use package context](https://www.youtube.com/watch?v=-_B5uQ4UGi0)
- [justforfunc #9: The Context Package](https://www.youtube.com/watch?v=LSzR0VEraWw)

| Name                                                      | Description                                        |
|-----------------------------------------------------------|----------------------------------------------------|
| [1-boring](/1-boring/main.go)                             | A hello world to goroutine                         |
| [2-chan](/2-chan/main.go)                                 | A hello world to go channel                        |
| [3-generator](/3-generator/main.go)                       | A python-liked generator                           |
| [4-fanin](/4-fanin/main.go)                               | Fan in pattern                                     |
| [5-restore-sequence](/5-restore-sequence/main.go)         | Restore sequence                                   |
| [6-select-timeout](/6-select-timeout/main.go)             | Add Timeout to a goroutine                         |
| [7-quit-signal](/7-quit-signal/main.go)                   | Quit signal                                        |
| [8-daisy-chan](/8-daisy-chan/main.go)                     | Daisy chan pattern                                 |
| [9-google1.0](/9-google1.0/main.go)                       | Build a concurrent google search from the grown-up |
| [10-google2.0](/10-google2.0/main.go)                     | Build a concurrent google search from the grown-up |
| [11-google2.1](/11-google2.1/main.go)                     | Build a concurrent google search from the grown-up |
| [12-google3.0](/12-google3.0/main.go)                     | Build a concurrent google search from the grown-up |
| [13-adv-pingpong](/13-adv-pingpong/main.go)               | A sample ping-pong table implemented in goroutine  |
| [14-adv-subscription](/14-adv-subscription/main.go)       | Subscription                                       |
| [15-bounded-parallelism](/15-bounded-parallelism/main.go) | Bounded parallelism                                |
| [16-context](/16-context/main.go)                         | How to user context in HTTP client and server      |
| [17-ring-buffer-channel](/17-ring-buffer-channel/main.go) | Ring buffer channel                                |
| [18-worker-pool](/18-worker-pool/main.go)                 | worker pool pattern                                |
