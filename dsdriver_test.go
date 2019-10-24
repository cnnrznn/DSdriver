package dsdriver

import (
	"testing"
)

func TestRemote(t *testing.T) {
	sChan, rChan := Remote(0)

	for i := 0; i < 10; i++ {
		sChan <- Message{[]byte("Hello!"), 0}
		<-rChan
	}
}
