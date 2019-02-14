package main

import (
	"bytes"
	"encoding/base64"
	_ "encoding/binary"
	"fmt"
	"io"
	"os"
	"time"
)

func BufferWriteRune() {
	var newRune = '好'
	buf := bytes.NewBufferString("Learning")
	fmt.Println(buf.String())

	buf.WriteRune(newRune)
	fmt.Println(buf.String())
}

func main() {
	fmt.Fprintf(os.Stdout, "os.Args:%v\n", len(os.Args))
	for key, arg := range os.Args {
		fmt.Print(key, "    ", arg)
		fmt.Print("\n")
	}
	fmt.Print("\n")

	str1 := bytes.NewBuffer([]byte("1234"))
	fmt.Print(str1, "\n")

	var b bytes.Buffer
	b.Write([]byte("bb"))
	b.WriteByte(byte('a'))
	//fmt.Print(b.Bytes());
	b.WriteString("Hello ")
	fmt.Fprintf(&b, "world! %s\n", "am niu")
	b.WriteTo(os.Stdout)

	buf := bytes.NewBufferString("R29waGVycyBydWxlIQ==")
	dec := base64.NewDecoder(base64.StdEncoding, buf)
	io.Copy(os.Stdout, dec)

	BufferWriteRune()

	fmt.Println("\n==============\n")
	var b2 = make([]byte, 10)
	fmt.Println("b2.len:", len(b2))
	var b3 = bytes.NewBuffer([]byte("234"))
	fmt.Println("b3.len:", b3.Len())

	select {
	//case m := <-time.After(2 * time.Second):
	//	fmt.Println("clock:", m)
	case m2 := <-time.NewTimer(1 * time.Second).C:
		fmt.Println("clock:", m2)
	}

	var arr_t [2]int = [2]int{1, 2}
	var slice_t []int = arr_t[:]
	fmt.Println("arr_t element:")
	for _, v := range arr_t {
		fmt.Print(v, " ")
	}
	fmt.Println("\nslice_t element:")
	for _, v := range slice_t {
		fmt.Print(v, " ")
	}
	var slice_t_2 []int

	slice_t_2 = make([]int, 5, 10)
	//slice_t_2[0] = 0
	fmt.Println("\nslice_t_2 element:")
	for _, v := range slice_t_2 {
		fmt.Print(v, " ")
	}
	fmt.Println()

}
