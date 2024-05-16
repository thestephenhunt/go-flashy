package utils

import (
	"math/rand"
	"strconv"
)

func atoi(s string) int {
	value, _ := strconv.Atoi(s)
	return value
}

func NewTerm() int {
	n := rand.Intn(100)
	return n
}

func Solve(f, s int, o string) int {
	switch o {
	case "+":
		return f + s
	default:
		return -1
	}
}

func ReverseString(str string) (result string) {
	for _, v := range str {
		result = string(v) + result
	}
	return
}
