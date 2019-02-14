package main

import (
	"fmt"
	"unsafe"
)

type slice struct {
	array unsafe.Pointer
	len   int
	cap   int
}

func main() {
	var a1 []int
	a1 = make([]int, 500)

	for i := 0; i < 500; i++ {
		a1[i] = i
	}
	fmt.Printf("a1: %v\n len:%d cap:%d\n", a1, len(a1), cap(a1))
	a1 = a1[:0]
	fmt.Printf("a1: %v\n len:%d cap:%d\n", a1, len(a1), cap(a1))
	a1 = append(a1, 898989)
	fmt.Printf("a1: %v\n len:%d cap:%d\n", a1, len(a1), cap(a1))
	return

	a2 := (*slice)(unsafe.Pointer(&a1))
	fmt.Printf("a2: %v\n", a2)
	fmt.Printf("a2.len:%d, a2.cap:%d\n", a2.len, a2.cap)
	a1 = append(a1, 8, 8, 9)
	fmt.Printf("a1: %v\n", a1)
	a3 := *(*[1]int)(a2.array)
	fmt.Printf("a3:%v\n", a3)

	fmt.Println()
	//a4 = append(a4, 8)
	a2.len = 3
	a1 = append(a1, 8)
	fmt.Printf("a1: %v\n", a1)

	a2 = (*slice)(unsafe.Pointer(&a1))
	fmt.Printf("a2: %v\n", a2)
	fmt.Printf("a2.len:%d, a2.cap:%d\n", a2.len, a2.cap)
	a3 = *(*[1]int)(a2.array)
	fmt.Printf("a3:%v\n", a3)

	return
}
