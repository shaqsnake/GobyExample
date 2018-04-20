package main

import (
	"fmt"
	"time"
)

func main() {
	requests := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		requests <- i
	}
	close(requests)

	limiter := time.Tick(time.Second)

	for req := range requests {
		<-limiter
		fmt.Println("request", req, time.Now())
	}

	brustyLimiter := make(chan time.Time, 3)

	for i := 0; i < 3; i++ {
		brustyLimiter <- time.Now()
	}

	go func() {
		for t := range time.Tick(time.Second) {
			brustyLimiter <- t
		}
	}()

	brustyRequests := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		brustyRequests <- i
	}
	close(brustyRequests)
	for req := range brustyRequests {
		<-brustyLimiter
		fmt.Println("request", req, time.Now())
	}
}
