package dsdriver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
    "net"
	"time"
)

type Dester interface {
	Dest() int // a function that reports the destination of this thing
}

type Message struct {
    Bytes []byte
    Dest int
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
	timeDiff := 50 * time.Millisecond

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

	go serve(i, nodes, toChan, frChan)

	return
}

func serve(i int, nodes []string, toChan, frChan chan Dester) {
    // TODO create listening socket
    pc, err := ListenPacket("udp", nodes[i])
    if err != nil {
        fmt.Println("Can't bind to socket", err)
        return
    }

}

func launchSenders(frChan chan Message) {
    for {
        select {
        case m := <-frChan:
            go send(m, nodes)
        }
    }
}

func send(data []byte, dest int, nodes []string) {
    for {
        conn, err := net.Dial("udp", nodes[dest])
        if err != nil {
            fmt.Println("Error dialing destination")
            continue
        }
        defer conn.Close()

        buf := make([]byte, 128)

        for {
            conn.Write(data)

            conn.SetReadDeadline(time.Now().Add(50 * time.Millisecond))

            n, err := conn.Read(buf)
            if err != nil {
                fmt.Println("Error reading packet")
                continue
            }
            if string(buf[:n]) == "ok" {
                return // Success!
            }
        }
    }
}
