package util

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/andreiz53/cookinator/types"
	"github.com/google/uuid"
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

// RandomMeasureUnit generates a random measure unit
func RandomMeasureUnit() types.MeasureUnit {
	return types.MeasureUnits[RandomInt(0, len(types.MeasureUnits)-1)]
}

// RandomRecipeItems generates an array of length 5 with recipe items
func RandomRecipeItems() []types.RecipeItem {
	var items []types.RecipeItem
	for i := 0; i < 5; i++ {
		items = append(items, types.RecipeItem{
			ID:       uuid.New(),
			Quantity: RandomFloat(1, 20),
			Unit:     RandomMeasureUnit(),
		})
	}

	return items
}
