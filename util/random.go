package util

import (
	"math/rand"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.New(rand.NewSource(time.Now().Unix()))
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var str string

	for i := 0; i < n; i++ {
		str += string(alphabet[rand.Intn(len(alphabet))])
	}
	return str
}

func RandomFirstName() string {
	return RandomString(8)
}

func RandomName() string {
	return RandomString(6)
}

func RandomEmail() string {
	return RandomString(5) + "@email.com"
}

func RandomPassword() string {
	return RandomString(16)
}
