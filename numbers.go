// Hoare 1978 4.5
// Set of integers in ascending order.
package main

import (
	"fmt"
	"time"
)

const nNums = 5

// num represents an integer which can print itself and accept a value.
// All input values come to the first integer which passes the value on to the next integer if necessary
// to preserve the sort order.
func num(iC chan int, nC chan int, qC chan bool, id int) {
	var content, i int
	empty := true
	lastN := false
	if id == nNums-1 {
		lastN = true
	}
	for {
		select {
		case <-qC:
			if empty {
				fmt.Printf("%v: empty\n", id)
				continue
			}
			fmt.Printf("%v: %v\n", id, content)
		case i = <-iC:
			switch {
			case empty:
				content = i
				empty = false
			case i < content:
				j := content
				content = i
				if !lastN {
					nC <- j
				}
			case i > content:
				if !lastN {
					nC <- i
				}
			}
		}
	}
}

func main() {
	var nC [nNums]chan int
	var qC [nNums]chan bool
	iC := make(chan int)
	for i := 0; i < nNums; i++ {
		nC[i] = make(chan int)
		qC[i] = make(chan bool)
		if i == 0 {
			go num(iC, nC[i], qC[i], i)
			continue
		}
		go num(nC[i-1], nC[i], qC[i], i) // The next channel of one integer becomes the input channel of the next.
	}

	for j := 0; j < nNums; j++ {
		qC[j] <- true
	}

	iC <- 900
	iC <- 300
	iC <- 800
	iC <- 850
	iC <- 200
	iC <- 600

	for j := 0; j < nNums; j++ {
		qC[j] <- true
	}

	iC <- 50
	iC <- 50

	for j := 0; j < nNums; j++ {
		qC[j] <- true
	}

	// main waits to see the output but then just exits without cleaning up.

	time.Sleep(10 * time.Second)
}
