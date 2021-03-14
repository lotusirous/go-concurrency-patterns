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

Context:
- [Go Concurrency Patterns: Context](https://blog.golang.org/context)
- [How to correctly use package context](https://www.youtube.com/watch?v=-_B5uQ4UGi0)
- [justforfunc #9: The Context Package](https://www.youtube.com/watch?v=LSzR0VEraWw)

| Name                   | Description                                        |
|------------------------|----------------------------------------------------|
| 1-boring               | A hello world to goroutine                         |
| 2-chan                 | A hello world to go channel                        |
| 3-generator            | A python-liked generator                           |
| 4-fanin                | Fan in pattern                                     |
| 5-restore-sequence     | Restore sequence                                   |
| 6-select-timeout       | Add Timeout to a goroutine                         |
| 7-quit-signal          | Quit signal                                        |
| 8-daisy-chan           | Daisy chan pattern                                 |
| 9-google1.0            | Build a concurrent google search from the grown-up |
| 10-google2.0           | Build a concurrent google search from the grown-up |
| 11-google2.1           | Build a concurrent google search from the grown-up |
| 12-google3.0           | Build a concurrent google search from the grown-up |
| 13-adv-pingpong        | A sample ping-pong table implemented in goroutine  |
| 14-adv-subscription    | Subcribtion                                        |
| 15-adv-subscription    | Bounded parallelism                                |
| 16-context             | How to user context in HTTP client and server      |
| 17-ring-buffer-channel | Ring buffer channel                                |
| 18-worker-pool         | worker pool pattern                                |
