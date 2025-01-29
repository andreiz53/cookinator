package util

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.New(rand.NewSource(time.Now().Unix()))
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

// RandomFloat generates a random float between min and max
func RandomFloat(min, max float64) float64 {
	return min + (rand.Float64() * (max - min))
}

// RandomPGNumeric generates a random pgtype.Numeric between 1 and 4
func RandomPGNumeric() pgtype.Numeric {
	var n pgtype.Numeric
	err := n.Scan(fmt.Sprint(RandomFloat(1, 4)))
	if err != nil {
		log.Fatal("cannot generate random pgtype.Numeric:", err)
	}
	return n
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var str string

	for i := 0; i < n; i++ {
		str += string(alphabet[rand.Intn(len(alphabet))])
	}
	return str
}

// RandomFirstName generates a random string of length 8
func RandomFirstName() string {
	return RandomString(8)
}

// RandomName generates a random string of length 6
func RandomName() string {
	return RandomString(6)
}

// RandomEmail generates a random email
func RandomEmail() string {
	return RandomString(5) + "@email.com"
}

// RandomPassword generates a random password of length 16
func RandomPassword() string {
	return RandomString(16)
}

func RandomTime() pgtype.Timestamp {
	var t pgtype.Timestamp
	randomTime := int(time.Now().Unix()) - RandomInt(1, 10000)
	err := t.Scan(fmt.Sprint(randomTime))
	if err != nil {
		log.Fatal("could not generate random time:", err)
	}
	return t
}
