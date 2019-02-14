package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	m := map[string]int{
		"abc": 123,
	}
	ms, _ := json.Marshal(m)
	fmt.Printf("m:%+v\nms:%+v\n", m, string(ms))

	b := make(map[string]int, 0)
	bs := `{"videoBaseName": "video_", "videoTableNum": 1,"commentBaseName": "comment","commentTableNum": 1 }`
	json.Unmarshal([]byte(bs), &b)

	fmt.Printf("b:%+v\nbs:%+v\n", b, string(bs))
	fmt.Println(b["commentBaseName"])
}
