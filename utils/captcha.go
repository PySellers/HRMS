package utils

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateCaptcha returns "answer|question"
func GenerateCaptcha() string {
	rand.Seed(time.Now().UnixNano())
	a := rand.Intn(9) + 1
	b := rand.Intn(9) + 1
	question := fmt.Sprintf("%d + %d = ?", a, b)
	answer := fmt.Sprintf("%d", a+b)
	return fmt.Sprintf("%s|%s", answer, question)
}
