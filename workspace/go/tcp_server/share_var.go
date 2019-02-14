package main

import "fmt"

import "time"

//共享变量有一个读通道和一个写通道组成
type shardedVar struct {
	reader chan int
	writer chan int
}

var value int = 0

//共享变量维护协程
func whachdog(v shardedVar) {
	go func() {
		//初始值
		for {
			//监听读写通道，完成服务
			select {
			case value = <-v.writer:
				print("write\n")
			case v.reader <- value:
				time.Sleep(1 * time.Second)
				print("read\n")
				value = 4
				//case <-time.After(1000 * time.Millisecond):
				//	fmt.Println("sleep .. value is:", value)
				//default:
				//	fmt.Println("default")
			}
		}
	}()
}
func main() {
	//初始化，并开始维护协程
	v := shardedVar{make(chan int), make(chan int)}
	whachdog(v)
	//读取初始值
	fmt.Println(<-v.reader)
	//写入一个值
	fmt.Println(<-v.reader)
	v.writer <- 1
	//读取新写入的值
	fmt.Println(<-v.reader)
	fmt.Println(<-v.reader)
}
