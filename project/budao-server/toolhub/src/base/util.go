package base

import (
	"math/rand"
	"time"
)

func ShuffleStringArr(array []string) {
	rand.Seed(time.Now().UTC().UnixNano())
	for i := len(array) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		array[i], array[j] = array[j], array[i]
	}
}
