package main

import (
	"fmt"
	"reflect"
)

type S1 struct {
	s1str string
	s1int int
}

func (s1 *S1) S1Func() {
	fmt.Println("this is S1Func")
}

//type S2 struct {
//	s1    *S1
//	s2str string
//	s2int int
//}
//
//func (s2 *S2) S2Func() (s1 *S1) {
//	fmt.Println("this is S2Func")
//	return s2.s1
//}

type S2 struct {
	s2str string
	s2int int
}

func (s2 *S2) S2Func() {
	fmt.Println("this is S2Func")
}

func main() {
	s2 := &S2{}
	//s2.s1 = &S1{}
	s2.s2str = "abc"
	s2.S2Func()
	t := s2
	fmt.Printf("t:%v\n", t)

	tType := reflect.TypeOf(t)
	fmt.Printf("tType:%v\n", tType)
	fmt.Printf("tType.Kind():%v\n", tType.Kind())
	for j := 0; j < tType.NumMethod(); j++ {
		fmt.Println(j, tType.Method(j).Name)
		j++
	}
	fmt.Printf("---------------------------\n")

	//fmt.Printf("tType.field[0]:%+v\n", tType.Field(0))

	tValue := reflect.ValueOf(t)
	fmt.Printf("tValue.Kind:%v\n", tValue.Kind())
	fmt.Printf("tValue.Elem:%v\n", tValue.Elem())
	fieldv := tValue.Elem().FieldByName(string("s2str")).CanInterface()
	fmt.Printf("----fieldv:%v\n", fieldv)
	fmt.Printf("tValue:%v\n", tValue)
	f := tValue.MethodByName("S2Func")
	nilValue := reflect.Value{}
	if f != nilValue {
		fmt.Printf("f:%+v\n", f)
		s1 := f.Call(nil)
		tValue1 := reflect.ValueOf(s1)
		f1 := tValue1.MethodByName("S1Func")
		fmt.Printf("tValue1:%v nilValue:%v f1:%v\n", tValue1, nilValue, f1)
		if f1 != nilValue {
			fmt.Printf("f1:%+v\n", f1)
			f1.Call(nil)

		}
	}

	fmt.Println("---------------------------")
	s := S2{}
	s.s2str = "abcd"
	sf := reflect.ValueOf(s).FieldByName("s2str")
	fmt.Printf("sf:%+v\n", sf)
}
