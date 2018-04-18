package main

import (
	"fmt"
)

func produce(ch chan int) {
	defer close(ch)

	for i := 0; i < 5; i++ {
		ch <- i
	}
}

func consume(ch chan int, done chan bool) {
	for v := range ch {
		fmt.Printf("%d ", v)
	}

	done <- true
}

func main() {
	ch := make(chan int)
	done := make(chan bool)

	go produce(ch)
	go consume(ch, done)

	// will block and wait done
	<-done
}
