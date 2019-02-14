package main

import (
	"bytes"
	"fmt"
	"net"
)

func ProcessBytes(buf []byte) {
	buf2 := bytes.Trim(buf, "\n")
	pos := bytes.Split(buf2, []byte(" "))
	fmt.Print("pos.len:", len(pos))
	i := 0
	for ; i < len(pos); i++ {
		fmt.Printf("pos[%d].len:%d\t", i, len(pos[i]))
		var str = string(pos[i])
		fmt.Printf("str.len:%d \t str:%s\n", len(str), str)
	}
	fmt.Printf("str:%c\n", pos[i-1][10])
}

func DoWork(conn net.Conn) {
	fmt.Println("new connection:", conn.LocalAddr())

	for {
		buf := make([]byte, 1024)
		length, err := conn.Read(buf)

		if err != nil {
			fmt.Println("Error reading:", err.Error())
			return
		}

		fmt.Println("Receive data from client:", string(buf[:length]), "=>length:", length)
		ProcessBytes(buf)
	}
}

func main() {
	fmt.Println("Start the server...")

	//create listener
	listener, err := net.Listen("tcp", "localhost:16379")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}

	//listen and accept connections from clients:
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
			return
		}
		//create a goroutine for each request.
		go DoWork(conn)
	}
}
