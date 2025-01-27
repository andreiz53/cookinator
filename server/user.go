package server

import (
	"net/http"
	"strings"

	database "github.com/andreiz53/cookinator/database/handlers"
	"github.com/andreiz53/cookinator/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type createUserRequest struct {
	FirstName string `json:"first_name" binding:"required,gte=3"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,gte=6"`
}

type createUserResponse struct {
	ID        uuid.UUID        `json:"id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
	FirstName string           `json:"first_name"`
	Email     string           `json:"email"`
	FamilyID  *uuid.UUID       `json:"family_id"`
}

func (s *Server) createUser(ctx *gin.Context) {
	var params createUserRequest

	err := ctx.ShouldBindJSON(&params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, respondWithErorr(err))
		return
	}

	dbParams := database.CreateUserParams{
		FirstName: params.FirstName,
		Email:     params.Email,
		Password:  params.Password,
	}

	user, err := s.store.CreateUser(ctx, dbParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, respondWithErorr(err))
		return
	}

	response := createUserResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		FirstName: user.FirstName,
		Email:     user.Email,
	}

	ctx.JSON(http.StatusCreated, response)
}

type getUsersResponse struct {
	ID        uuid.UUID        `json:"id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
	FirstName string           `json:"first_name"`
	Email     string           `json:"email"`
	FamilyID  *uuid.UUID       `json:"family_id"`
}

func (s *Server) getUsers(ctx *gin.Context) {
	users, err := s.store.GetUsers(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, respondWithErorr(err))
		return
	}

	var response []getUsersResponse
	for _, user := range users {
		response = append(response, getUsersResponse{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			FirstName: user.FirstName,
			Email:     user.Email,
			FamilyID:  util.NullUUID(user.FamilyID),
		})
	}

	ctx.JSON(http.StatusOK, response)
}

type getUserByIDRequest struct {
	ID string `uri:"id" binding:"required,uuid4_rfc4122"`
}

type getUserByIDResponse struct {
	ID        uuid.UUID        `json:"id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
	FirstName string           `json:"first_name"`
	Email     string           `json:"email"`
	FamilyID  *uuid.UUID       `json:"family_id"`
}

func (s *Server) getUserByID(ctx *gin.Context) {
	var request getUserByIDRequest

	err := ctx.ShouldBindUri(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, respondWithErorr(err))
		return
	}

	user, err := s.store.GetUserByID(ctx, uuid.MustParse(request.ID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, respondWithErorr(err))
		return
	}

	response := getUserByIDResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		FirstName: user.FirstName,
		Email:     user.Email,
		FamilyID:  util.NullUUID(user.FamilyID),
	}

	ctx.JSON(http.StatusOK, response)
}

type updateUserEmailRequest struct {
	ID    string `json:"id" binding:"required,uuid4_rfc4122"`
	Email string `json:"email" binding:"required,email"`
}

type updateUserEmailResponse struct {
	ID        uuid.UUID        `json:"id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
	FirstName string           `json:"first_name"`
	Email     string           `json:"email"`
	FamilyID  *uuid.UUID       `json:"family_id"`
}

func (s *Server) updateUserEmail(ctx *gin.Context) {
	var request updateUserEmailRequest

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, respondWithErorr(err))
		return
	}

	dbParams := database.UpdateUserEmailParams{
		ID:    uuid.MustParse(request.ID),
		Email: request.Email,
	}
	user, err := s.store.UpdateUserEmail(ctx, dbParams)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			ctx.JSON(http.StatusConflict, respondWithErorr(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, respondWithErorr(err))
		return
	}

	response := updateUserEmailResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		FirstName: user.FirstName,
		Email:     user.Email,
		FamilyID:  util.NullUUID(user.FamilyID),
	}

	ctx.JSON(http.StatusOK, response)
}

func (s *Server) updateUserPassword(ctx *gin.Context) {}
func (s *Server) updateUserInfo(ctx *gin.Context)     {}
