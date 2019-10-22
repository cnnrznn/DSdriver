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
	Dest() int               // the message's destination
	Encode() ([]byte, error) // serialize the message
	Decode([]byte) error     // de-serialize the message
}

type Hub func(sendChans []chan Message, recvChan chan Message)

func BenignHub(sendChans []chan Message, recvChan chan Message) {
	for {
		select {
		case d := <-recvChan:
			sendChans[d.Dest()] <- d
		}
	}
}

func ReorderHub(sendChans []chan Message, recvChan chan Message) {
	buffer := make([]Message, 0)
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

func Local(n int, fn Hub) (frChan chan Message,
	toChans []chan Message) {
	frChan = make(chan Message, 1024)

	for i := 0; i < n; i++ {
		toChans = append(toChans, make(chan Message, 1024))
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

func Remote(i int) (frChan, toChan chan Message) {
	nodes, err := loadNodes()
	if err != nil {
		fmt.Println("Error loading 'nodes' file", err)
		panic("Error loading file")
	}
	fmt.Println(nodes)

	frChan = make(chan Message, 1024)
	toChan = make(chan Message, 1024)

	go serve(i, nodes, toChan)

	for {
		select {
		case msg := <-frChan:
			go send(msg, nodes)
		}
	}

	return
}

func serve(i int, nodes []string, toChan chan Message) {
}

func send(msg Message, nodes []string) {
}
