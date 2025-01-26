package database

import (
	"context"
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/andreiz53/cookinator/types"
	"github.com/andreiz53/cookinator/util"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func createRandomRecipe(t *testing.T) Recipe {
	family := createRandomFamily(t)
	recipeItems := util.RandomRecipeItems()
	recipeItemsData, err := json.Marshal(recipeItems)
	if err != nil {
		log.Fatal("could not stringify json recipe items:", err)
	}
	arg := CreateRecipeParams{
		Name:           util.RandomName(),
		CookingProcess: util.RandomString(128),
		FamilyID:       family.ID,
		Items:          recipeItemsData,
	}

	recipe, err := testQueries.CreateRecipe(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, recipe)

	require.Equal(t, arg.Name, recipe.Name)
	require.Equal(t, arg.CookingProcess, recipe.CookingProcess)
	require.Equal(t, arg.FamilyID, recipe.FamilyID)
	require.NotZero(t, recipe.ID)

	checkRecipeItems(t, recipeItemsData, recipe.Items)

	return recipe
}

func checkRecipeItems(t *testing.T, items1 []byte, items2 []byte) {
	var recipe1Items []types.RecipeItem
	var recipe2Items []types.RecipeItem

	err := json.Unmarshal(items1, &recipe1Items)
	require.NoError(t, err)

	err = json.Unmarshal(items2, &recipe2Items)
	require.NoError(t, err)

	require.Equal(t, recipe1Items, recipe2Items)
}

func TestCreateRecipe(t *testing.T) {
	createRandomRecipe(t)
}

func TestGetRecipeByID(t *testing.T) {
	recipe := createRandomRecipe(t)

	recipe2, err := testQueries.GetRecipeByID(context.Background(), recipe.ID)
	require.NoError(t, err)
	require.NotEmpty(t, recipe2)

	require.Equal(t, recipe.ID, recipe2.ID)
	require.Equal(t, recipe.Name, recipe2.Name)
	require.Equal(t, recipe.CookingProcess, recipe2.CookingProcess)
	require.Equal(t, recipe.FamilyID, recipe2.FamilyID)

	require.WithinDuration(t, recipe.CreatedAt.Time, recipe2.CreatedAt.Time, time.Second)

	checkRecipeItems(t, recipe.Items, recipe2.Items)
}

func TestGetRecipes(t *testing.T) {
	createRandomRecipe(t)
	createRandomRecipe(t)
	createRandomRecipe(t)

	recipes, err := testQueries.GetRecipes(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, recipes)

	require.True(t, len(recipes) >= 3)
}

func TestGetRecipesByFamilyID(t *testing.T) {
	recipe := createRandomRecipe(t)

	recipes, err := testQueries.GetRecipesByFamilyID(context.Background(), recipe.FamilyID)
	require.NoError(t, err)
	require.NotEmpty(t, recipes)

	require.True(t, len(recipes) >= 1)
}

func TestUpdateRecipe(t *testing.T) {
	recipe := createRandomRecipe(t)
	recipeItems := util.RandomRecipeItems()
	recipeItemsData, err := json.Marshal(recipeItems)
	if err != nil {
		log.Fatal("could not stringify json recipe items:", err)
	}

	arg := UpdateRecipeParams{
		ID:             recipe.ID,
		Name:           util.RandomName(),
		CookingProcess: util.RandomString(128),
		Items:          recipeItemsData,
	}

	recipe2, err := testQueries.UpdateRecipe(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, recipe2)

	require.Equal(t, recipe.ID, recipe2.ID)
	require.Equal(t, arg.Name, recipe2.Name)
	require.Equal(t, arg.CookingProcess, recipe2.CookingProcess)
	require.Equal(t, recipe.FamilyID, recipe2.FamilyID)

	require.WithinDuration(t, recipe.CreatedAt.Time, recipe2.CreatedAt.Time, time.Second)

	checkRecipeItems(t, recipeItemsData, recipe2.Items)
}

func TestDeleteRecipe(t *testing.T) {
	recipe := createRandomRecipe(t)

	err := testQueries.DeleteRecipe(context.Background(), recipe.ID)
	require.NoError(t, err)

	recipe2, err := testQueries.GetRecipeByID(context.Background(), recipe.ID)
	require.Error(t, err)
	require.Empty(t, recipe2)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
}
