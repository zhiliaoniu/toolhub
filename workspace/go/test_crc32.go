package main

import (
	"fmt"
	"hash/crc32"
	"math/rand"
	"time"
)

func main() {
	str := make([]byte, 0)
	str = append(str, string("abcdefghijklmnopqrstuvwxyz1234567890")...)
	l := len(str)
	result := make(map[uint32]uint32, 0)
	randSource := rand.NewSource(time.Now().UnixNano())
	var Max uint32 = 10000000
	var Num uint32 = 10
	var i uint32 = 0
	for ; i < Max; i++ {
		random := rand.New(randSource).Uint32() % uint32(l)
		str[0], str[random] = str[random], str[0]
		crc := crc32.ChecksumIEEE([]byte(str))
		result[crc%Num] += 1
	}
	for _, count := range result {
		fmt.Printf("count:%d percent:%8.5f\n", count, float64(count)/float64(Max))
	}
}
