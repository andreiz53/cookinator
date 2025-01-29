// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package database

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateFamily(ctx context.Context, arg CreateFamilyParams) (Family, error)
	CreateIngredient(ctx context.Context, arg CreateIngredientParams) (Ingredient, error)
	CreateRecipe(ctx context.Context, arg CreateRecipeParams) (Recipe, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteFamily(ctx context.Context, id uuid.UUID) error
	DeleteIngredient(ctx context.Context, id int32) error
	DeleteRecipe(ctx context.Context, id uuid.UUID) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	GetFamilies(ctx context.Context) ([]Family, error)
	GetFamilyByID(ctx context.Context, id uuid.UUID) (Family, error)
	GetFamilyByUserID(ctx context.Context, createdByUserID uuid.UUID) (Family, error)
	GetIngredientByID(ctx context.Context, id int32) (Ingredient, error)
	GetIngredientByName(ctx context.Context, name string) (Ingredient, error)
	GetIngredients(ctx context.Context) ([]Ingredient, error)
	GetRecipeByID(ctx context.Context, id uuid.UUID) (Recipe, error)
	GetRecipes(ctx context.Context) ([]Recipe, error)
	GetRecipesByFamilyID(ctx context.Context, familyID uuid.UUID) ([]Recipe, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (User, error)
	GetUsers(ctx context.Context) ([]User, error)
	UpdateFamily(ctx context.Context, arg UpdateFamilyParams) (Family, error)
	UpdateIngredient(ctx context.Context, arg UpdateIngredientParams) (Ingredient, error)
	UpdateRecipe(ctx context.Context, arg UpdateRecipeParams) (Recipe, error)
	UpdateUserEmail(ctx context.Context, arg UpdateUserEmailParams) (User, error)
	UpdateUserInfo(ctx context.Context, arg UpdateUserInfoParams) (User, error)
	UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) (User, error)
}

var _ Querier = (*Queries)(nil)
