package dsdriver

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
	"math/rand"
	"time"
)

type Dester interface {
	Dest() int // a function that reports the destination of this thing
}

type Message interface {
    Dest() int // the message's destination
    Encode() []byte, error // serialize the message
    Decode([]byte) error // de-serialize the message
}

type Hub func(sendChans []chan Dester, recvChan chan Dester)

func BenignHub(sendChans []chan Dester, recvChan chan Dester) {
	for {
		select {
		case d := <-recvChan:
			sendChans[d.Dest()] <- d
		}
	}
}

func ReorderHub(sendChans []chan Dester, recvChan chan Dester) {
	buffer := make([]Dester, 0)
	startTime := time.Now()
	timeDiff := 100 * time.Millisecond

	for {
		select {
		case d := <-recvChan:
			buffer = append(buffer, d)
		default:
			if time.Now().Sub(startTime) > timeDiff {
				startTime = time.Now()
				rand.Shuffle(len(buffer), func(i, j int) {
					buffer[i], buffer[j] = buffer[j], buffer[i]
				})
				for len(buffer) > 0 {
					d := buffer[0]
					buffer = buffer[1:]
					sendChans[d.Dest()] <- d
				}
			}
		}
	}
}

func Local(n int, fn Hub) (frChan chan Dester,
	toChans []chan Dester) {
	frChan = make(chan Dester, 1024)

	for i := 0; i < n; i++ {
		toChans = append(toChans, make(chan Dester, 1024))
	}

	go fn(toChans, frChan)

	return
}

func loadNodes() (nodes []string, err error) {
    data, err := ioutil.ReadFile("nodes.json")
    if err != nil {
        return
    }

    err = json.Unmarshal(data, &nodes)

    return
}

func Remote(i int) (frChan, toChan chan Dester) {
    nodes, err := loadNodes()
    if err != nil {
        fmt.Println("Error loading 'nodes' file", err)
        panic("Error loading file")
    }
    fmt.Println(nodes)

    frChan = make(chan Dester, 1024)
    toChan = make(chan Dester, 1024)

    go serve(i, nodes, toChan)

    for {
        select {
        case msg := <-frChan:
            go send(msg, nodes)
        }
    }

    return
}
