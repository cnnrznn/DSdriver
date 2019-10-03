package dsdriver

import (
    "container/list"
    "fmt"
    "math/rand"
    "sync"
    "time"
)

type replica interface {
    New() replica       // need a new instance of the protocol
    Run()               // a function that runs the protocol
}

type event interface {
    Dest() int
}

func hub(sendChans []chan event,
         recvChan chan event,
         delay int) {
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

func System(replicas []replica) {
    toreps := []chan event{}
    fromreps := make(chan event, 1024)

    var wg sync.WaitGroup
    wg.Add(len(replicas))
    defer wg.Wait()

    fmt.Print("Initialization...")

    rand.Seed(time.Now().UnixNano())

    for i := 0; i < len(replicas); i++ {
        toreps = append(toreps, make(chan event, 1024))
    }

    fmt.Println("Done.")

    for _, r := range replicas {
        go r.Run()
    }

    go hub(toreps, fromreps, 0)
}

