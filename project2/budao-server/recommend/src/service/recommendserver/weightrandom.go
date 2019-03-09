package recommendserver

import (
	"crypto/rand"
	"errors"
	"math/big"
)

func IntRange(min, max int) (int, error) {
	var result int
	switch {
	case min > max:
		return result, errors.New("Min cannot be greater than max.")
	case max == min:
		result = max
	case max > min:
		maxRand := max - min
		b, err := rand.Int(rand.Reader, big.NewInt(int64(maxRand)))
		if err != nil {
			return result, err
		}
		result = min + int(b.Int64())
	}
	return result, nil
}

type Choice struct {
	Weight int
	Item   interface{}
}

func WeightedRandom(choices []Choice) (Choice, error) {
	var ret Choice
	sum := 0
	for _, c := range choices {
		sum += c.Weight
	}
	r, err := IntRange(0, sum)
	if err != nil {
		return ret, err
	}
	for _, c := range choices {
		r -= c.Weight
		if r < 0 {
			return c, nil
		}
	}
	err = errors.New("Internal error - code should not reach this point")
	return ret, err
}
