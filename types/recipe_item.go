package types

import (
	"github.com/google/uuid"
)

type RecipeItem struct {
	ID       uuid.UUID   `json:"id"`
	Quantity float64     `json:"quantity"`
	Unit     MeasureUnit `json:"unit"`
}
