// Hoare 1978 5.1
// Bounded buffer.
package main

import (
	"fmt"
	"math/rand"
	"time"
)

const lotLen = 10

// The buffer consists of 10 lots each a 10 element array.
func buffer(inC chan [lotLen]int, outC chan [lotLen]int) {
	const numLots = 10
	var buf [numLots][lotLen]int
	var inCt, outCt int
	inCOK := true
	allDone := false
	for {
		switch {
		case outCt < inCt:
			outC <- buf[outCt%numLots]
			outCt++
		case !inCOK && inCt == outCt:
			allDone = true // The termination condition is in = out and producer channel is closed.
		case inCt < outCt+numLots:
			buf[inCt%numLots], inCOK = <-inC
			if inCOK {
				inCt++
			}
		}
		if allDone {
			break
		}
	}
	close(outC)
}

// The producer just fires off 100 lots with some random waits.
func produce(outC chan [lotLen]int) {
	var lot [lotLen]int
	var c int
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		for j := 0; j < lotLen; j++ {
			lot[j] = c
			c++
		}
		time.Sleep(time.Duration(rand.Intn(15)) * time.Millisecond)
		outC <- lot
	}
	close(outC)
}

// The consumer is held up for a random amount of time as it processes each lot.
func consume(inC chan [lotLen]int, ctrlC chan bool) {
	for c := range inC {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		fmt.Println(c)
	}
	ctrlC <- true // signal main
}

func main() {
	iC := make(chan [lotLen]int)
	oC := make(chan [lotLen]int)
	cC := make(chan bool)
	go buffer(iC, oC)
	go produce(iC)
	go consume(oC, cC)
	<-cC // Wait for consumer process to finish.
}
