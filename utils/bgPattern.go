package utils

import "math/rand"

var availableBgs = [6]string{"wavy-bg", "red-square", "pink-green-swirl", "color-squares", "circle-halves", "circle-braid"}
var BgPattern = func() string {
	return availableBgs[rand.Intn(len(availableBgs))]
}
