package main

import (
	"fmt"
	"os"
	"bytes"
	"encoding/base64"
	"io"
	_"encoding/binary"
)

func BufferWriteRune() {
	var newRune = '好'
	buf := bytes.NewBufferString("Learning");
	fmt.Println(buf.String())
	
	buf.WriteRune(newRune)
	fmt.Println(buf.String())
}

func main () {
	str1 := bytes.NewBuffer([]byte("1234"));
	fmt.Print(str1, "\n")
	

	var b bytes.Buffer
	b.Write([]byte("bb"))
	b.WriteByte(byte('a'))
	//fmt.Print(b.Bytes());
	b.WriteString("Hello ")
	fmt.Fprintf(&b, "world!\n")
	b.WriteTo(os.Stdout)

	buf := bytes.NewBufferString("R29waGVycyBydWxlIQ==")
	dec := base64.NewDecoder(base64.StdEncoding, buf)
	io.Copy(os.Stdout, dec)

	BufferWriteRune()
}
