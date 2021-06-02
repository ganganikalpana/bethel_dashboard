package utils

import (
	"math/rand"
	"time"
)

func GenerateCode() int {
	rand.Seed(time.Now().UnixNano())
	rn := rand.Intn(100000)
	return rn
}
