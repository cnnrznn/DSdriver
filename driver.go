package main

import (
    "container/list"
    "fmt"
    "github.com/cnnrznn/rcft"
    "math/rand"
    "os"
    "strconv"
    "sync"
    "time"
)

func hub(sendChans []chan rcft.Event,
         recvChan chan rcft.Event,
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
                e := buffer.Front().Value.(rcft.Event)
                sendChans[e.Pid] <- e
                buffer.Remove(buffer.Front())
            }
        }
    }
}

func instance(n, f, delay int) {
    replicas := []rcft.Replica{}
    toreps := []chan rcft.Event{}
    fromreps := make(chan rcft.Event, 1024)

    var wg sync.WaitGroup
    wg.Add(n)
    defer wg.Wait()

    fmt.Print("Initialization...")

    rand.Seed(time.Now().UnixNano())

    for i := 0; i < n; i++ {
        replicas = append(replicas, rcft.NewReplica(rand.Intn(2)))
        toreps = append(toreps, make(chan rcft.Event, 1024))
    }

    fmt.Println("Done.")

    for i, r := range replicas {
        go r.Consensus(n, f, fromreps, toreps[i], &wg)
    }

    go hub(toreps, fromreps, delay)
}

func main() {
    args := os.Args[1:]
    n, _ := strconv.Atoi(args[0])
    f, _ := strconv.Atoi(args[1])
    delay, _ := strconv.Atoi(args[2])

    instance(n, f, delay)
}

