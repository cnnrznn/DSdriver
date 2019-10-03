package dsdriver

import (
    "container/list"
    "fmt"
    "math/rand"
    "sync"
    "time"
)

type replica interface {
    Run(n, f int, fr, to chan event) // a function that runs the protocol
}

type event interface {
    Dest() int          // a function that reports the destination of the message
}

func hub(sendChans []chan event, recvChan chan event, delay int) {
    buffer := list.New()

    for {
        select {
        case e := <-recvChan:
            fmt.Println(e)
            buffer.PushBack(e)
        default:
            time.Sleep(time.Duration(delay) * time.Millisecond)
            for buffer.Len() > 0 {
                e := buffer.Front().Value.(event)
                sendChans[e.Dest()] <- e
                buffer.Remove(buffer.Front())
            }
        }
    }
}

func run(r replica, n, f int, fr, to chan event, wg *sync.WaitGroup) {
    r.Run(n, f, fr, to)
    wg.Done()
}

func System(replicas []replica, n, f int) {
    toreps := []chan event{}
    fromreps := make(chan event, 1024)

    var wg sync.WaitGroup
    wg.Add(len(replicas))
    defer wg.Wait()

    fmt.Print("Initialization...")

    rand.Seed(time.Now().UnixNano())

    for i := 0; i < n; i++ {
        toreps = append(toreps, make(chan event, 1024))
    }

    fmt.Println("Done.")

    for i, r := range replicas {
        go run(r, n, f, fromreps, toreps[i], &wg)
    }

    go hub(toreps, fromreps, 0)
}

