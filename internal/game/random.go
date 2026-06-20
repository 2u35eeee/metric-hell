package game

import (
	"math/rand"
	"time"
)

func NewRandomSeed() int64 {
	return time.Now().UnixNano() + rand.Int63n(997)
}
