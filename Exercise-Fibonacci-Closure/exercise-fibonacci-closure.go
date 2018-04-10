package main

import "fmt"

func fibonacci() func() int {
	x1, x2 := 1, 0
	return func() int {
		x1, x2 = x2, x1+x2
		return x1
	}
}

func main() {
	f := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
}
