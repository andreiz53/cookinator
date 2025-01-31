package server

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	database "github.com/andreiz53/cookinator/database/handlers"
	databaseMock "github.com/andreiz53/cookinator/database/mocks"
	"github.com/andreiz53/cookinator/util"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateIngredient(t *testing.T) {
	ingredientParams := database.CreateIngredientParams{
		Name:    util.RandomName(),
		Density: util.RandomPGNumeric(),
	}

	testCases := []struct {
		name          string
		params        database.CreateIngredientParams
		stubs         func(store *databaseMock.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			params: ingredientParams,
			stubs: func(store *databaseMock.MockStore) {
				store.EXPECT().
					CreateIngredient(mock.Anything, ingredientParams).
					Times(1).
					Return(randomIngredient(), nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name:   "BadRequest",
			params: database.CreateIngredientParams{},
			stubs: func(store *databaseMock.MockStore) {
				store.EXPECT().
					CreateIngredient(mock.Anything, database.CreateIngredientParams{}).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "Duplicated",
			params: ingredientParams,
			stubs: func(store *databaseMock.MockStore) {
				store.EXPECT().
					CreateIngredient(mock.Anything, mock.Anything).
					Times(1).
					Return(database.Ingredient{}, database.ErrDuplicateKey)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusConflict, recorder.Code)
			},
		},
		{
			name:   "InternalServerError",
			params: ingredientParams,
			stubs: func(store *databaseMock.MockStore) {
				store.EXPECT().
					CreateIngredient(mock.Anything, mock.Anything).
					Times(1).
					Return(database.Ingredient{}, pgx.ErrTxClosed)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			store := new(databaseMock.MockStore)
			server := NewServer(store)

			tc.stubs(store)

			recorder := httptest.NewRecorder()
			url := "/ingredients"

			data, err := encodeJSON(tc.params)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetIngredients(t *testing.T) {
	var ingredients []database.Ingredient
	n := 3
	for i := 0; i < n; i++ {
		ingredients = append(ingredients, randomIngredient())
	}

	testCases := []struct {
		name          string
		ingredients   []database.Ingredient
		stubs         func(store *databaseMock.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:        "OK",
			ingredients: ingredients,
			stubs: func(store *databaseMock.MockStore) {
				store.EXPECT().
					GetIngredients(mock.Anything).
					Times(1).
					Return(ingredients, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchIngredients(t, recorder.Body, ingredients)
			},
		},
		{
			name:        "InternalServerError",
			ingredients: []database.Ingredient{},
			stubs: func(store *databaseMock.MockStore) {
				store.EXPECT().
					GetIngredients(mock.Anything).
					Times(1).
					Return([]database.Ingredient{}, pgx.ErrTxClosed)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			store := new(databaseMock.MockStore)
			server := NewServer(store)

			tc.stubs(store)

			recorder := httptest.NewRecorder()
			url := "/ingredients"

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetIngredientByID(t *testing.T) {
	ingredient := randomIngredient()

	testCases := []struct {
		name          string
		ingredientID  int32
		stubs         func(store *databaseMock.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:         "OK",
			ingredientID: ingredient.ID,
			stubs: func(store *databaseMock.MockStore) {
				store.EXPECT().
					GetIngredientByID(mock.Anything, ingredient.ID).
					Times(1).
					Return(ingredient, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchIngredient(t, recorder.Body, ingredient)
			},
		},
		{
			name:         "NotFound",
			ingredientID: ingredient.ID,
			stubs: func(store *databaseMock.MockStore) {
				store.EXPECT().
					GetIngredientByID(mock.Anything, ingredient.ID).
					Times(1).
					Return(database.Ingredient{}, pgx.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:         "InternalServerError",
			ingredientID: ingredient.ID,
			stubs: func(store *databaseMock.MockStore) {
				store.EXPECT().
					GetIngredientByID(mock.Anything, ingredient.ID).
					Times(1).
					Return(database.Ingredient{}, pgx.ErrTxClosed)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:         "BadRequest",
			ingredientID: 0,
			stubs: func(store *databaseMock.MockStore) {
				store.EXPECT().
					GetIngredientByID(mock.Anything, mock.Anything).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			store := new(databaseMock.MockStore)
			server := NewServer(store)

			tc.stubs(store)

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/ingredients/%d", tc.ingredientID)

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}

}

func TestUpdateIngredient(t *testing.T) {}

func TestDeleteIngredient(t *testing.T) {}

func randomIngredient() database.Ingredient {
	return database.Ingredient{
		ID:      int32(util.RandomInt(1, 1000)),
		Name:    util.RandomName(),
		Density: util.RandomPGNumeric(),
	}
}

// func createRandomIngredient(t *testing.T) Ingredient {
// 	params := database.CreateIngredientParams{
// 		Name:    util.RandomName(),
// 		Density: util.RandomPGNumeric(),
// 	}

// 	expectedValue := database.Ingredient{
// 		ID:      1,
// 		Name:    params.Name,
// 		Density: params.Density,
// 	}

// 	store := new(databaseMock.MockStore)
// 	server := NewServer(store)
// 	gin.SetMode(gin.TestMode)
// 	store.EXPECT().CreateIngredient(mock.Anything, params).Times(1).Return(expectedValue, nil)

// 	recorder := httptest.NewRecorder()
// 	url := "/ingredients"

// 	data, err := encodeJSON(params)
// 	require.NoError(t, err)
// 	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
// 	require.NoError(t, err)
// 	request.Header.Set("Content-Type", "application/json")

// 	server.router.ServeHTTP(recorder, request)

// 	require.Equal(t, http.StatusCreated, recorder.Code)
// 	return requireBodyMatchIngredient(t, recorder.Body, expectedValue)
// }

func requireBodyMatchIngredient(t *testing.T, body *bytes.Buffer, ingredient database.Ingredient) {
	createdIngredient, err := decodeJSON[Ingredient](body)
	require.NoError(t, err)
	require.NotEmpty(t, createdIngredient)

	require.Equal(t, ingredient.Name, createdIngredient.Name)
	require.Equal(t, ingredient.Density, createdIngredient.Density)
}

func requireBodyMatchIngredients(t *testing.T, body *bytes.Buffer, ingredients []database.Ingredient) {
	createdIngredients, err := decodeJSON[[]Ingredient](body)
	require.NoError(t, err)
	require.NotEmpty(t, createdIngredients)

	require.Equal(t, len(ingredients), len(createdIngredients))
	for i, ingredient := range createdIngredients {
		require.Equal(t, ingredient.Name, ingredients[i].Name)
		require.Equal(t, ingredient.Density, ingredients[i].Density)
		require.Equal(t, ingredient.ID, ingredients[i].ID)
	}
}
