package server

import (
	"net/http"
	"strings"

	database "github.com/andreiz53/cookinator/database/handlers"
	"github.com/andreiz53/cookinator/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) createUser(ctx *gin.Context) {
	var params types.CreateUserParams

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

	response := types.User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		FirstName: user.FirstName,
		Email:     user.Email,
	}

	ctx.JSON(http.StatusCreated, response)
}

func (s *Server) getUsers(ctx *gin.Context) {
	users, err := s.store.GetUsers(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, respondWithErorr(err))
		return
	}

	ctx.JSON(http.StatusOK, types.DBUsersToUsers(users))
}

func (s *Server) getUserByID(ctx *gin.Context) {
	var request types.GetUserByIDParams

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

	ctx.JSON(http.StatusOK, types.DBUserToUser(user))
}

func (s *Server) updateUserEmail(ctx *gin.Context) {
	var request types.UpdateUserEmailParams

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

	ctx.JSON(http.StatusOK, types.DBUserToUser(user))
}

func (s *Server) updateUserPassword(ctx *gin.Context) {}
func (s *Server) updateUserInfo(ctx *gin.Context)     {}
