package utils

import (
	"math"
	"math/rand"
	"time"
)

func Init() {
	rand.Seed(time.Now().UnixNano())
}

func MaxI64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func MaxUI32(a, b uint32) uint32 {
	if a > b {
		return a
	}
	return b
}

func MaxF64(a, b float64) float64 {
	return math.Max(a, b)
}

func RandIntInRange(min, max int) uint32 {
	bid := uint32(min + rand.Intn(max-min))
	return bid
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
