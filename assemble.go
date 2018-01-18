// From Hoare 1978 3.4
// Read 40 character "cards" and print 65 character "lines".
package main

import "fmt"

func assemble(in chan rune, ctrl chan bool) {
	const lnLen = 65
	line := make([]rune, lnLen)
	var lnIdx int
	for b := range in {
		line[lnIdx] = b
		lnIdx++
		if lnIdx == lnLen {
			fmt.Printf("%c\n", line)
			lnIdx = 0
		}
	}
	if lnIdx > 0 {
		for i := lnIdx; i < lnLen; i++ {
			line[i] = ' '
		}
		fmt.Printf("%c\n", line)
	}
	ctrl <- true
}

func main() {
	var cards [3]string
	cards[0] = "0123456789012345678901234567890123456789"
	cards[1] = "0123456789012345678901234567890123456789"
	cards[2] = "0123456789012345678901234567890123456789"
	bChan := make(chan rune)
	cChan := make(chan bool)
	go assemble(bChan, cChan)
	for _, c := range cards {
		for _, b := range c {
			bChan <- b
		}
	}
	close(bChan)
	<-cChan
	close(cChan)
}
