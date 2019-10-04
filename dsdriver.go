package dsdriver

import (
    "container/list"
    "time"
)

type Dester interface {
    Dest() int          // a function that reports the destination of this thing
}

func hub(sendChans []chan Dester, recvChan chan Dester, delay int) {
    buffer := list.New()

    for {
        select {
        case d := <-recvChan:
            buffer.PushBack(d)
        default:
            time.Sleep(time.Duration(delay) * time.Millisecond)
            for buffer.Len() > 0 {
                d := buffer.Front().Value.(Dester)
                sendChans[d.Dest()] <- d
                buffer.Remove(buffer.Front())
            }
        }
    }
}

func Local(n int) (frChan chan Dester,
                    toChans []chan Dester) {
    frChan = make(chan Dester, 1024)

    for i := 0; i < n; i++ {
        toChans = append(toChans, make(chan Dester, 1024))
    }

    go hub(toChans, frChan, 0)

    return
}

