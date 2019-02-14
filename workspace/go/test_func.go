package main

import (
	"fmt"
	"reflect"
)

type Test struct {
	testMem  string
	testMem2 string
}

func (s *Test) Func() (r string) {
	fmt.Printf("call Func\n")
	return "result from Func"
}
func (s *Test) Func2() {
	fmt.Printf("call Func2\n")
}

func main() {
	t := &Test{}
	fmt.Printf("%+v\n", reflect.ValueOf(t))
	v := reflect.ValueOf(t)
	r := v.MethodByName("Func").Call(nil)
	fmt.Printf("result.len:%d\n", len(r))
	fmt.Printf("result:%+v\n", r[0])

	var s string
	s = r[0].String()
	fmt.Printf("s:%s\n", string(s))
	fmt.Printf("exist:%+v\n", reflect.ValueOf(t).MethodByName("Func"))
	fmt.Printf("not exist:%+v\n", reflect.ValueOf(t).MethodByName("Func3"))
	null := reflect.Value{}
	if reflect.ValueOf(t).MethodByName("Func3") == null {
		fmt.Printf("ok")
	}
}
