package types

import (
	database "github.com/andreiz53/cookinator/database/handlers"
	"github.com/andreiz53/cookinator/util"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID        uuid.UUID        `json:"id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
	FirstName string           `json:"first_name"`
	Email     string           `json:"email"`
	FamilyID  *uuid.UUID       `json:"family_id"`
}

type CreateUserParams struct {
	FirstName string `json:"first_name" binding:"required,min=3"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
}

type GetUserByIDParams struct {
	ID string `uri:"id" binding:"required,uuid4_rfc4122"`
}

type UpdateUserEmailParams struct {
	ID    string `json:"id" binding:"required,uuid4_rfc4122"`
	Email string `json:"email" binding:"required,email"`
}

type UpdateUserPasswordParams struct {
	ID       string `json:"id" binding:"required,uuid4_rfc4122"`
	Password string `json:"password" binding:"required,min=6"`
}

type UpdateUserInfoParams struct {
	ID        string `json:"id" binding:"required,uuid4_rfc4122"`
	FirstName string `json:"first_name" binding:"required,min=3"`
}

type DeleteUserParams struct {
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
