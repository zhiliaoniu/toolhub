package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func main() {
	cmd := exec.Command("ls", "-ht")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))
	arr := strings.Split(string(out), "\n")
	fmt.Printf("arr:%v\n", arr)
	fmt.Printf("arr[0]:%v\n", arr[0])
}
