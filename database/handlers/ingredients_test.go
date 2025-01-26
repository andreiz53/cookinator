package database

import (
	"context"
	"testing"

	"github.com/andreiz53/cookinator/util"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func createRandomIngredient(t *testing.T) Ingredient {
	arg := CreateIngredientParams{
		Name:    util.RandomName(),
		Density: util.RandomPGNumeric(),
	}

	ingredient, err := testQueries.CreateIngredient(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, ingredient)

	require.Equal(t, arg.Name, ingredient.Name)
	require.Equal(t, arg.Density, ingredient.Density)
	require.NotZero(t, ingredient.ID)

	return ingredient
}

func TestCreateIngredient(t *testing.T) {
	createRandomIngredient(t)
}

func TestGetIngredientByID(t *testing.T) {
	ingredient := createRandomIngredient(t)

	ingredient2, err := testQueries.GetIngredientByID(context.TODO(), ingredient.ID)
	require.NoError(t, err)
	require.NotEmpty(t, ingredient2)

	require.Equal(t, ingredient.ID, ingredient2.ID)
	require.Equal(t, ingredient.Name, ingredient2.Name)
	require.Equal(t, ingredient.Density, ingredient2.Density)
}

func TestGetIngredientByName(t *testing.T) {
	ingredient := createRandomIngredient(t)

	ingredient2, err := testQueries.GetIngredientByName(context.TODO(), ingredient.Name)
	require.NoError(t, err)
	require.NotEmpty(t, ingredient2)

	require.Equal(t, ingredient.ID, ingredient2.ID)
	require.Equal(t, ingredient.Name, ingredient2.Name)
	require.Equal(t, ingredient.Density, ingredient2.Density)
}

func TestGetIngredients(t *testing.T) {
	createRandomIngredient(t)
	createRandomIngredient(t)
	createRandomIngredient(t)

	ingredients, err := testQueries.GetIngredients(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, ingredients)

	require.True(t, len(ingredients) >= 3)
}

func TestUpdateIngredient(t *testing.T) {
	ingredient := createRandomIngredient(t)

	arg := UpdateIngredientParams{
		ID:      ingredient.ID,
		Name:    util.RandomName(),
		Density: util.RandomPGNumeric(),
	}

	ingredient2, err := testQueries.UpdateIngredient(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, ingredient2)

	require.Equal(t, ingredient.ID, ingredient2.ID)
	require.Equal(t, arg.Name, ingredient2.Name)
	require.Equal(t, arg.Density, ingredient2.Density)
}

func TestDeleteIngredient(t *testing.T) {
	ingredient := createRandomIngredient(t)

	err := testQueries.DeleteIngredient(context.Background(), ingredient.ID)
	require.NoError(t, err)

	ingredient2, err := testQueries.GetIngredientByID(context.Background(), ingredient.ID)
	require.Error(t, err)
	require.Empty(t, ingredient2)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
}
