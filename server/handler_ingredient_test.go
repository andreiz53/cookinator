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
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func randomIngredient() database.Ingredient {
	return database.Ingredient{
		ID:      int32(util.RandomInt(1, 1000)),
		Name:    util.RandomName(),
		Density: util.RandomPGNumeric(),
	}
}

func createRandomIngredient(t *testing.T) Ingredient {
	params := database.CreateIngredientParams{
		Name:    util.RandomName(),
		Density: util.RandomPGNumeric(),
	}

	expectedValue := database.Ingredient{
		ID:      1,
		Name:    params.Name,
		Density: params.Density,
	}

	store := new(databaseMock.MockStore)
	server := NewServer(store)
	gin.SetMode(gin.TestMode)
	store.EXPECT().CreateIngredient(mock.Anything, params).Times(1).Return(expectedValue, nil)

	recorder := httptest.NewRecorder()
	url := "/ingredients"

	data, err := encodeJSON(params)
	require.NoError(t, err)
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	require.NoError(t, err)
	request.Header.Set("Content-Type", "application/json")

	server.router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusCreated, recorder.Code)
	return requireBodyMatchIngredient(t, recorder.Body, expectedValue)
}

func requireBodyMatchIngredient(t *testing.T, body *bytes.Buffer, ingredient database.Ingredient) Ingredient {
	createdIngredient, err := decodeJSON[Ingredient](body)
	require.NoError(t, err)
	require.NotEmpty(t, createdIngredient)

	require.Equal(t, ingredient.Name, createdIngredient.Name)
	require.Equal(t, ingredient.Density, createdIngredient.Density)

	return createdIngredient
}

func TestCreateIngredient(t *testing.T) {
	createRandomIngredient(t)
}

func TestGetIngredients(t *testing.T) {

}

func TestGetIngredientByID(t *testing.T) {
	ingredient := database.Ingredient{
		ID:      int32(util.RandomInt(100, 10000)),
		Name:    util.RandomName(),
		Density: util.RandomPGNumeric(),
	}

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
