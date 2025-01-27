package database

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// RecipeItem corresponds to the PostgreSQL `recipe_item` type
type RecipeItem struct {
	ID       uuid.UUID      `json:"id"`
	Quantity pgtype.Numeric `json:"quantity"`
	Unit     string         `json:"unit"`
}

type MeasureUnit string

const (
	MeasureUnitGrams       = "g"
	MeasureUnitMillilitres = "mL"
	MeasureUnitTeaspoon    = "tsp"
	MeasureUnitTablespoon  = "tbsp"
	MeasureUnitPiece       = "pc"
	MeasureUnitCup         = "cup"
)

var MeasureUnits = []MeasureUnit{
	MeasureUnitGrams,
	MeasureUnitCup,
	MeasureUnitMillilitres,
	MeasureUnitPiece,
	MeasureUnitTablespoon,
	MeasureUnitTeaspoon,
}
