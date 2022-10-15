# Scheduler

A tool for scheduling the execution of a function at a given interval.

## Example

```go
package main

import (
    "context"
    "fmt"
    "github.com/maxoov1/scheduler"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func sayHello(name string) {
    fmt.Printf("Hello, %s! \n", name)
}

func main() {
    sc := scheduler.NewScheduler()

    opts := scheduler.ExecuteJobOptions{
        Job:       sayHello,
        Arguments: []any{"Maksim"},
        Timeout:   1 * time.Second,
    }

    if err := sc.ExecuteJob(context.Background(), opts); err != nil {
        panic(err)
    }

    exit := make(chan os.Signal, 1)
    signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

    <-exit

    sc.Shutdown()
}
```
