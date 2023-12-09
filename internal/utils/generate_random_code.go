package utils

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateRandomCode generates a random 4-digit code as a string.
func GenerateRandomCode() string {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	return fmt.Sprintf("%04d", r.Intn(10000))
}
