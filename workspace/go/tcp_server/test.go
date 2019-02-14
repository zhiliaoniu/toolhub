package main

import (
	_ "bytes"
	"fmt"
	"math/rand"
	"runtime"
	_ "strings"
)

func randGenerator() chan int {
	ch := make(chan int)
	go func() {
		ch <- rand.Int()
	}()
	return ch
}

func doSomething(num int) (sum int) {
	for i := 1; i <= 10; i++ {
		fmt.Printf("%d + %d = %d\n", num, num+i, num+num+i)
		sum = sum + num + i
	}
	return sum
}
func testLoop() {
	// 建立计数器，通道大小为cpu核数
	var NumCPU = runtime.NumCPU()
	fmt.Printf("NumCPU = %d\n", NumCPU)
	sem := make(chan int, NumCPU)
	//FOR循环体
	data := []int{1, 11, 21, 31, 41, 51, 61, 71, 81, 91}
	for _, v := range data {
		//建立协程
		go func(v int) {
			fmt.Printf("doSomething(%d)...\n", v)
			sum := doSomething(v)
			//计数
			sem <- sum
		}(v)
	}
	// 等待循环结束
	var total int = 0
	for i := 0; i < len(data); i++ {
		temp := <-sem
		fmt.Printf("%d <- sem\n", temp)
		total = total + temp
	}
	fmt.Printf("total = %d\n", total)
}

//func main() {
//	str := "aaayangbbb"
//	strl := strings.Trim(str, "ab")
//	fmt.Printf("str:%s\tstrl:%s\n", str, strl)
//	var buf bytes.Buffer
//	buf.WriteString("aaayangbbb")
//	fmt.Printf("b:%s b.len:%d\n", buf.Bytes(), len(buf.Bytes()))
//
//	ch := randGenerator()
//	rand := <-ch
//	fmt.Println("rand:", rand)
//
//	testLoop()
//}
func Generate(ch chan int) {
	for i := 2; ; i++ {
		print("Generate beforei:", i, "\n")
		ch <- i // Send 'i' to channel 'in'.
		print("Generate i:", i, "\n")
	}
}
func Filter(in chan int, out chan int, prime int) {
	for {
		print("Filter out before.prime:", prime, "\n")
		i := <-in // Receive valuefrom 'in'.
		print("Filter out prime:", prime, "\t i:", i, "\n")
		if i%prime != 0 {
			print("Filter in prime:", prime, "\t i:", i, "\n")
			out <- i // Send'i' to 'out'.
			print("Filter in after.prime:", prime, "\t i:", i, "\n")
		}
	}
}
func main() {
	in := make(chan int)
	go Generate(in)

	for i := 0; i < 3; i++ {
		print("main before in:", i, "\n")
		prime := <-in
		print(prime, "\n")
		out := make(chan int)
		go Filter(in, out, prime)
		in = out
	}
}
