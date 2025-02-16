package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	database "github.com/andreiz53/cookinator/database/handlers"
	"github.com/andreiz53/cookinator/util"
)

type createUserRequest struct {
	FirstName string `json:"first_name" binding:"required,min=3"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
}

type User struct {
	ID        uuid.UUID        `json:"id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
	FirstName string           `json:"first_name"`
	Email     string           `json:"email"`
	FamilyID  *uuid.UUID       `json:"family_id"`
}

type getUserByIDRequest struct {
	ID string `uri:"id" binding:"required,uuid4_rfc4122"`
}

type updateUserEmailRequest struct {
	ID    string `json:"id" binding:"required,uuid4_rfc4122"`
	Email string `json:"email" binding:"required,email"`
}

type updateUserPasswordRequest struct {
	ID       string `json:"id" binding:"required,uuid4_rfc4122"`
	Password string `json:"password" binding:"required,min=6"`
}

type updateUserInfoRequest struct {
	ID        string `json:"id" binding:"required,uuid4_rfc4122"`
	FirstName string `json:"first_name" binding:"required,min=3"`
}

type deleteUserRequest struct {
	ID string `uri:"id" binding:"required,uuid4_rfc4122"`
}

func DBUserToUser(arg database.User) User {
	return User{
		ID:        arg.ID,
		CreatedAt: arg.CreatedAt,
		UpdatedAt: arg.UpdatedAt,
		FirstName: arg.FirstName,
		Email:     arg.Email,
		FamilyID:  util.NullUUID(arg.FamilyID),
	}
}

func DBUsersToUsers(arg []database.User) []User {
	users := []User{}
	for _, user := range arg {
		users = append(users, DBUserToUser(user))
	}
	return users
}

func (s *Server) createUser(ctx *gin.Context) {
	var request createUserRequest

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, respondWithErorr(err))
		return
	}
	hashedPassword, err := util.HashPassword(request.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, respondWithErorr(err))
		return
	}

	dbParams := database.CreateUserParams{
		FirstName: request.FirstName,
		Email:     request.Email,
		Password:  hashedPassword,
	}

	user, err := s.store.CreateUser(ctx, dbParams)
	if err != nil {
		if database.ErrorCode(err) == database.CodeDuplicateKey {
			ctx.JSON(http.StatusConflict, respondWithErorr(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, respondWithErorr(err))
		return
	}

	ctx.JSON(http.StatusCreated, DBUserToUser(user))
}

func (s *Server) getUsers(ctx *gin.Context) {
	users, err := s.store.GetUsers(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, respondWithErorr(err))
		return
	}

	ctx.JSON(http.StatusOK, DBUsersToUsers(users))
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

	ctx.JSON(http.StatusOK, DBUserToUser(user))
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
		if database.ErrorCode(err) == database.CodeDuplicateKey {
			ctx.JSON(http.StatusConflict, respondWithErorr(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, respondWithErorr(err))
		return
	}

	ctx.JSON(http.StatusOK, DBUserToUser(user))
}

func (s *Server) updateUserPassword(ctx *gin.Context) {
	var request updateUserPasswordRequest

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, respondWithErorr(err))
		return
	}

	hashedPassword, err := util.HashPassword(request.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, respondWithErorr(err))
		return
	}

	dbParams := database.UpdateUserPasswordParams{
		ID:       uuid.MustParse(request.ID),
		Password: hashedPassword,
	}
	user, err := s.store.UpdateUserPassword(ctx, dbParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, respondWithErorr(err))
		return
	}

	ctx.JSON(http.StatusOK, DBUserToUser(user))
}

func (s *Server) updateUserInfo(ctx *gin.Context) {
	var request updateUserInfoRequest

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, respondWithErorr(err))
		return
	}

	dbParams := database.UpdateUserInfoParams{
		ID:        uuid.MustParse(request.ID),
		FirstName: request.FirstName,
	}
	user, err := s.store.UpdateUserInfo(ctx, dbParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, respondWithErorr(err))
		return
	}

	ctx.JSON(http.StatusOK, DBUserToUser(user))
}

func (s *Server) deleteUser(ctx *gin.Context) {
	var request deleteUserRequest

	err := ctx.ShouldBindUri(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, respondWithErorr(err))
		return
	}

	err = s.store.DeleteUser(ctx, uuid.MustParse(request.ID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, respondWithErorr(err))
		return
	}

	ctx.JSON(http.StatusOK, respondWithMessage(fmt.Sprintf("deleted user with id %s", request.ID)))
}
