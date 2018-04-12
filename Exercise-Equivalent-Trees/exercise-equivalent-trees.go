package main

import (
	"fmt"

	"golang.org/x/tour/tree"
)

func Walk(t *tree.Tree, ch chan int) {
	_walk(t, ch)
	close(ch)
}

func _walk(t *tree.Tree, ch chan int) {
	if t == nil {
		return
	}
	_walk(t.Left, ch)
	ch <- t.Value
	_walk(t.Right, ch)
}

func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go Walk(t1, ch1)
	go Walk(t2, ch2)

	var v1, v2 int
	var ok bool
	for v1 = range ch1 {
		v2, ok = <-ch2
		if !ok || v1 != v2 {
			return false
		}
	}
	_, ok = <-ch2
	return !ok
}

func main() {
	ch := make(chan int)
	t1 := tree.New(1)
	go Walk(t1, ch)
	for v := range ch {
		fmt.Printf("%d ", v)
	}
	fmt.Print("\n")

	ch = make(chan int)
	t2 := tree.New(2)
	go Walk(t2, ch)
	for v := range ch {
		fmt.Printf("%d ", v)
	}
	fmt.Print("\n")

	fmt.Println(Same(t1, t2))
	fmt.Println(Same(t1, t1))
}
