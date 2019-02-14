package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func DoRead(conn net.Conn) {
	//do read
	defer func() {
		fmt.Println("got err\n")
		if err := recover(); err != nil {
			fmt.Println("err:", err)
		}
	}()
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		panic(err.Error())
	}
}

func main() {
	//open conntion:
	conn, err := net.Dial("tcp", "localhost:16379")
	if err != nil {
		fmt.Println("Error dial:", err.Error())
		return
	}

	inputReader := bufio.NewReader(os.Stdin)
	fmt.Println("Please input your name:")
	clientName, _ := inputReader.ReadString('\n')
	inputClientName := strings.Trim(clientName, "\n")

	go DoRead(conn)

	//send info to server until Quit
	for {
		fmt.Println("What do you want send to the server? Type Q to quit.")
		content, _ := inputReader.ReadString('\n')
		inputContent := strings.Trim(content, "\n")
		if inputContent == "Q" {
			return
		}

		_, err := conn.Write([]byte(inputClientName + " say: " + inputContent))
		if err != nil {
			fmt.Println("Error Write:", err.Error())
			return
		}
	}
}
