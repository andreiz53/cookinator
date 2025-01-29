package types

import (
	"fmt"

	database "github.com/andreiz53/cookinator/database/handlers"
	"github.com/jackc/pgx/v5/pgtype"
)

type Ingredient struct {
	ID      int32          `json:"id"`
	Name    string         `json:"name"`
	Density pgtype.Numeric `json:"density"`
}

type CreateIngredientParams struct {
	Name    string  `json:"name" binding:"required,min=2"`
	Density float64 `json:"density" binding:"required,gt=0"`
}

type UpdateIngredientParams struct {
	ID      int32   `json:"id" binding:"required,number"`
	Name    string  `json:"name" binding:"required,min=2"`
	Density float64 `json:"density" binding:"required,gt=0"`
}

type DeleteIngredientParams struct {
	ID int32 `uri:"id" binding:"required,number"`
}

type GetIngredientByIDParams struct {
	ID int32 `uri:"id" binding:"required,number"`
}

func CreateIngredientToDBCreateIngredient(arg CreateIngredientParams) database.CreateIngredientParams {
	var density pgtype.Numeric
	err := density.Scan(fmt.Sprint(arg.Density))
	if err != nil {
		density.Valid = false
		density.NaN = true
	}
	return database.CreateIngredientParams{
		Name:    arg.Name,
		Density: density,
	}
}

func UpdateIngredientToDBUpdateIngredient(arg UpdateIngredientParams) database.UpdateIngredientParams {
	var density pgtype.Numeric
	err := density.Scan(fmt.Sprint(arg.Density))
	if err != nil {
		density.Valid = false
		density.NaN = true
	}
	return database.UpdateIngredientParams{
		Name:    arg.Name,
		Density: density,
		ID:      arg.ID,
	}
}

func DBIngredientToIngredient(arg database.Ingredient) Ingredient {
	return Ingredient{
		ID:      arg.ID,
		Name:    arg.Name,
		Density: arg.Density,
	}
}

func DBIngredientsToIngredients(arg []database.Ingredient) []Ingredient {
	var ingredients []Ingredient
	for _, ing := range arg {
		ingredients = append(ingredients, DBIngredientToIngredient(ing))
	}
	return ingredients
}
