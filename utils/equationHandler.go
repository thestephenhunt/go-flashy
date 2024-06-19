package utils

import (
	"log"
	"math/rand"
	"strconv"
	"strings"
)

func NewTerm(d int) int {
	maxNum, _ := strconv.Atoi(strings.Repeat("9", d))
	log.Println(maxNum)
	n := rand.Intn(maxNum)
	return n
}

func Solve(f, s int, o string) int {
	switch o {
	case "+":
		return f + s
	case "-":
		return f - s
	case "*":
		return f * s
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
