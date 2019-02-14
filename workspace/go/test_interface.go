package main

import (
	"fmt"
	"strings"
	"time"
)

var Max int = 100000

func TransStrArrToInterface(strArr []string) []interface{} {
	l := len(strArr)
	retArr := make([]interface{}, l)

	//retArr = append(retArr, []interface{}(strArr)...)
	for i := 0; i < l; i++ {
		retArr[i] = strArr[i]
	}
	return retArr
}

func TransStrArrToInterface2(arr []interface{}) []interface{} {
	retArr := make([]interface{}, 0)
	retArr = append(retArr, arr...)
	fmt.Printf("arr:%v\nretArr:%v\n", arr, retArr)
	return retArr
}

func main() {
	beginT := time.Now()
	/*
		strArr := ""
		for n := 0; n < 5000; n++ {
			strArr += "123456789abc,"
		}
		for n := 0; n < Max; n++ {
			_ = strings.Split(strArr, ",")
		}
	*/
	strArr := make([]string, 0)
	for n := 0; n < 5000; n++ {
		strArr = append(strArr, "123456789abc")
	}

	for i := 0; i < Max; i++ {
		//TransStrArrToInterface(strArr)
		_ = strings.Join(strArr, ",")
	}
	fmt.Printf("num:%d, cost time:%d", Max, time.Now().Sub(beginT))
}
