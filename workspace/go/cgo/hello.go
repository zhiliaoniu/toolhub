package main

import "fmt"

/*
#cgo  CFLAGS:  -I  ./include
#cgo  LDFLAGS:  -L ./lib  -lhello -lstdc++
#include "hello.h"
*/
import "C"

func main() {
	fmt.Println("vim-go")
	C.hello(C.CString("haha"))
}
