package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/andreiz53/cookinator/types"
	"github.com/gin-gonic/gin"
)

func (s *Server) createIngredient(ctx *gin.Context) {
	var request types.CreateIngredientParams

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, respondWithErorr(err))
		return
	}

	ingredient, err := s.store.CreateIngredient(ctx, types.CreateIngredientToDBCreateIngredient(request))
	if err != nil {
		ctx.JSON(http.StatusConflict, respondWithErorr(err))
		return
	}

	ctx.JSON(http.StatusCreated, types.DBIngredientToIngredient(ingredient))
}

func (s *Server) getIngredients(ctx *gin.Context) {
	ingredients, err := s.store.GetIngredients(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, respondWithErorr(err))
		return
	}
	ctx.JSON(http.StatusOK, types.DBIngredientsToIngredients(ingredients))
}

func (s *Server) getIngredientByID(ctx *gin.Context) {
	var request types.GetIngredientByIDParams

	err := ctx.ShouldBindUri(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, respondWithErorr(err))
		return
	}

	ingredient, err := s.store.GetIngredientByID(ctx, request.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, respondWithErorr(err))
		return
	}

	ctx.JSON(http.StatusOK, types.DBIngredientToIngredient(ingredient))
}

func (s *Server) updateIngredient(ctx *gin.Context) {
	var request types.UpdateIngredientParams

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, respondWithErorr(err))
		return
	}

	ingredient, err := s.store.UpdateIngredient(ctx, types.UpdateIngredientToDBUpdateIngredient(request))
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			ctx.JSON(http.StatusConflict, respondWithErorr(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, respondWithErorr(err))
		return
	}

	ctx.JSON(http.StatusOK, types.DBIngredientToIngredient(ingredient))
}

func (s *Server) deleteIngredient(ctx *gin.Context) {
	var request types.DeleteIngredientParams

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
