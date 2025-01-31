package server

import (
	"fmt"
	"net/http"
	"strings"

	database "github.com/andreiz53/cookinator/database/handlers"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
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
	ID      int32   `json:"id" binding:"required,min=1"`
	Name    string  `json:"name" binding:"required,min=2"`
	Density float64 `json:"density" binding:"required,gt=0"`
}

type DeleteIngredientParams struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

type GetIngredientByIDParams struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func createIngredientToDBCreateIngredient(arg CreateIngredientParams) database.CreateIngredientParams {
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

func updateIngredientToDBUpdateIngredient(arg UpdateIngredientParams) database.UpdateIngredientParams {
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

func dbIngredientToIngredient(arg database.Ingredient) Ingredient {
	return Ingredient{
		ID:      arg.ID,
		Name:    arg.Name,
		Density: arg.Density,
	}
}

func dbIngredientsToIngredients(arg []database.Ingredient) []Ingredient {
	var ingredients []Ingredient
	for _, ing := range arg {
		ingredients = append(ingredients, dbIngredientToIngredient(ing))
	}
	return ingredients
}

func (s *Server) createIngredient(ctx *gin.Context) {
	var request CreateIngredientParams

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, respondWithErorr(err))
		return
	}

	ingredient, err := s.store.CreateIngredient(ctx, createIngredientToDBCreateIngredient(request))
	if err != nil {
		ctx.JSON(http.StatusConflict, respondWithErorr(err))
		return
	}
	ctx.JSON(http.StatusCreated, dbIngredientToIngredient(ingredient))
}

func (s *Server) getIngredients(ctx *gin.Context) {
	ingredients, err := s.store.GetIngredients(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, respondWithErorr(err))
		return
	}
	ctx.JSON(http.StatusOK, dbIngredientsToIngredients(ingredients))
}

func (s *Server) getIngredientByID(ctx *gin.Context) {
	var request GetIngredientByIDParams

	err := ctx.ShouldBindUri(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, respondWithErorr(err))
		return
	}

	ingredient, err := s.store.GetIngredientByID(ctx, request.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, respondWithErorr(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, respondWithErorr(err))
		return
	}
	ctx.JSON(http.StatusOK, dbIngredientToIngredient(ingredient))
}

func (s *Server) updateIngredient(ctx *gin.Context) {
	var request UpdateIngredientParams

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, respondWithErorr(err))
		return
	}

	ingredient, err := s.store.UpdateIngredient(ctx, updateIngredientToDBUpdateIngredient(request))
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			ctx.JSON(http.StatusConflict, respondWithErorr(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, respondWithErorr(err))
		return
	}

	ctx.JSON(http.StatusOK, dbIngredientToIngredient(ingredient))
}

func (s *Server) deleteIngredient(ctx *gin.Context) {
	var request DeleteIngredientParams

	err := ctx.ShouldBindUri(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, respondWithErorr(err))
		return
	}

	err = s.store.DeleteIngredient(ctx, request.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, respondWithErorr(err))
		return
	}
	ctx.JSON(http.StatusOK, respondWithMessage(fmt.Sprintf("deleted ingredient with id %d", request.ID)))
}
