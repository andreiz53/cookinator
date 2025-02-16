package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	database "github.com/andreiz53/cookinator/database/handlers"
)

type Family struct {
	ID              uuid.UUID        `json:"id"`
	CreatedAt       pgtype.Timestamp `json:"created_at"`
	UpdatedAt       pgtype.Timestamp `json:"updated_at"`
	Name            string           `json:"name"`
	CreatedByUserID uuid.UUID        `json:"created_by_user_id"`
}

type CreateFamilyParams struct {
	Name   string `json:"name" binding:"required,min=2"`
	UserID string `json:"user_id" binding:"required,uuid4_rfc4122"`
}

type GetFamilyByIDParams struct {
	ID string `uri:"id" binding:"required,uuid4_rfc4122"`
}

type GetFamilyByUserIDParams struct {
	UserID string `uri:"user_id" binding:"required,uuid4_rfc4122"`
}

type UpdateFamilyParams struct {
	ID   string `json:"id" binding:"required,uuid4_rfc4122"`
	Name string `json:"name" binding:"required,min=2"`
}

type DeleteFamilyParams struct {
	ID string `uri:"id" binding:"required,uuid4_rfc4122"`
}

func createFamilyToDBCreateFamily(arg CreateFamilyParams) database.CreateFamilyParams {
	return database.CreateFamilyParams{
		Name:            arg.Name,
		CreatedByUserID: uuid.MustParse(arg.UserID),
	}
}

func DBFamilyToFamily(arg database.Family) Family {
	return Family{
		ID:              arg.ID,
		CreatedAt:       arg.CreatedAt,
		UpdatedAt:       arg.UpdatedAt,
		Name:            arg.Name,
		CreatedByUserID: arg.CreatedByUserID,
	}
}

func DBFamiliesToFamilies(arg []database.Family) []Family {
	var families []Family
	for _, family := range arg {
		families = append(families, DBFamilyToFamily(family))
	}
	return families
}

func UpdateFamilyToDBUpdateFamily(arg UpdateFamilyParams) database.UpdateFamilyParams {
	return database.UpdateFamilyParams{
		ID:   uuid.MustParse(arg.ID),
		Name: arg.Name,
	}
}

func (s *Server) createFamily(ctx *gin.Context) {
	var request CreateFamilyParams
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, respondWithErorr(err))
		return
	}

	family, err := s.store.CreateFamily(ctx, createFamilyToDBCreateFamily(request))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, respondWithErorr(err))
		return
	}
	ctx.JSON(http.StatusCreated, DBFamilyToFamily(family))
}

func (s *Server) getFamilies(ctx *gin.Context) {
	families, err := s.store.GetFamilies(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, respondWithErorr(err))
		return
	}

	ctx.JSON(http.StatusOK, DBFamiliesToFamilies(families))
}

func (s *Server) getFamilyByID(ctx *gin.Context) {
	var request GetFamilyByIDParams

	err := ctx.ShouldBindUri(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, respondWithErorr(err))
		return
	}

	family, err := s.store.GetFamilyByID(ctx, uuid.MustParse(request.ID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, respondWithErorr(err))
		return
	}

	ctx.JSON(http.StatusOK, DBFamilyToFamily(family))
}

func (s *Server) getFamilyByUserID(ctx *gin.Context) {
	var request GetFamilyByUserIDParams

	err := ctx.ShouldBindUri(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, respondWithErorr(err))
		return
	}

	family, err := s.store.GetFamilyByUserID(ctx, uuid.MustParse(request.UserID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, respondWithErorr(err))
		return
	}

	ctx.JSON(http.StatusOK, DBFamilyToFamily(family))
}

func (s *Server) updateFamily(ctx *gin.Context) {
	var request UpdateFamilyParams

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, respondWithErorr(err))
		return
	}

	family, err := s.store.UpdateFamily(ctx, UpdateFamilyToDBUpdateFamily(request))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, respondWithErorr(err))
		return
	}

	ctx.JSON(http.StatusOK, DBFamilyToFamily(family))
}

func (s *Server) deleteFamily(ctx *gin.Context) {
	var request DeleteFamilyParams

	err := ctx.ShouldBindUri(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, respondWithErorr(err))
		return
	}

	err = s.store.DeleteFamily(ctx, uuid.MustParse(request.ID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, respondWithErorr(err))
		return
	}

	ctx.JSON(http.StatusOK, respondWithMessage(fmt.Sprintf("deleted family with id %s", request.ID)))
}
