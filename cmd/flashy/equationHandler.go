package main

import (
	"math/rand"
)

func NewTerm() int64 {
	n := rand.Int63n(100)
	return n
}
