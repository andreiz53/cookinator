package server

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	database "github.com/andreiz53/cookinator/database/handlers"
	databaseMock "github.com/andreiz53/cookinator/database/mocks"
	"github.com/andreiz53/cookinator/util"
)

func randomFamily() database.Family {
	return database.Family{
		ID:              uuid.New(),
		CreatedAt:       util.RandomTime(),
		UpdatedAt:       util.RandomTime(),
		Name:            util.RandomName(),
		CreatedByUserID: uuid.New(),
	}
}

func TestGetFamilies(t *testing.T) {
	var families []database.Family
	n := 3
	for i := 0; i < n; i++ {
		families = append(families, randomFamily())
	}

	testCases := []struct {
		name          string
		families      []database.Family
		stubs         func(store *databaseMock.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			families: families,
			stubs: func(store *databaseMock.MockStore) {
				store.EXPECT().
					GetFamilies(mock.Anything).
					Times(1).Return(families, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:     "InternalServerError",
			families: families,
			stubs: func(store *databaseMock.MockStore) {
				store.EXPECT().
					GetFamilies(mock.Anything).
					Times(1).Return([]database.Family{}, pgx.ErrTxClosed)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			store := new(databaseMock.MockStore)
			server := newTestServer(t, store)

			tc.stubs(store)

			recorder := httptest.NewRecorder()
			url := "/families"

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetFamilyByID(t *testing.T) {
	family := randomFamily()

	testCases := []struct {
		name          string
		familyID      uuid.UUID
		stubs         func(store *databaseMock.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			familyID: family.ID,
			stubs: func(store *databaseMock.MockStore) {
				store.EXPECT().
					GetFamilyByID(mock.Anything, family.ID).
					Times(1).Return(family, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:     "BadRequest",
			familyID: uuid.UUID{},
			stubs: func(store *databaseMock.MockStore) {
				store.EXPECT().
					GetFamilyByID(mock.Anything, family.ID).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:     "NotFound",
			familyID: family.ID,
			stubs: func(store *databaseMock.MockStore) {
				store.EXPECT().
					GetFamilyByID(mock.Anything, family.ID).
					Times(1).Return(database.Family{}, pgx.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			store := new(databaseMock.MockStore)
			server := newTestServer(t, store)

			tc.stubs(store)

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/families/%s", tc.familyID.String())

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetFamilyByUserID(t *testing.T) {
	family := randomFamily()

	testCases := []struct {
		name          string
		userID        uuid.UUID
		stubs         func(store *databaseMock.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			userID: family.CreatedByUserID,
			stubs: func(store *databaseMock.MockStore) {
				store.EXPECT().
					GetFamilyByUserID(mock.Anything, family.CreatedByUserID).
					Times(1).Return(family, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "BadRequest",
			userID: uuid.UUID{},
			stubs: func(store *databaseMock.MockStore) {
				store.EXPECT().
					GetFamilyByUserID(mock.Anything, family.CreatedByUserID).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "NotFound",
			userID: family.CreatedByUserID,
			stubs: func(store *databaseMock.MockStore) {
				store.EXPECT().
					GetFamilyByUserID(mock.Anything, family.CreatedByUserID).
					Times(1).Return(database.Family{}, pgx.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			store := new(databaseMock.MockStore)
			server := newTestServer(t, store)

			tc.stubs(store)

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/families/users/%s", tc.userID.String())

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestCreateFamily(t *testing.T) {
	family := randomFamily()

	params := CreateFamilyParams{
		Name:   util.RandomName(),
		UserID: family.CreatedByUserID.String(),
	}

	testCases := []struct {
		name          string
		params        CreateFamilyParams
		stubs         func(store *databaseMock.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			params: params,
			stubs: func(store *databaseMock.MockStore) {
				store.EXPECT().
					CreateFamily(mock.Anything, createFamilyToDBCreateFamily(params)).
					Times(1).Return(family, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name:   "BadRequest",
			params: CreateFamilyParams{},
			stubs: func(store *databaseMock.MockStore) {
				store.EXPECT().
					CreateFamily(mock.Anything, mock.Anything).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "InternalServerError",
			params: params,
			stubs: func(store *databaseMock.MockStore) {
				store.EXPECT().
					CreateFamily(mock.Anything, createFamilyToDBCreateFamily(params)).
					Times(1).Return(database.Family{}, pgx.ErrTxClosed)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			store := new(databaseMock.MockStore)
			server := newTestServer(t, store)

			tc.stubs(store)

			recorder := httptest.NewRecorder()
			url := "/families"

			data, err := encodeJSON(tc.params)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestUpdateFamily(t *testing.T) {
	family := randomFamily()

	params := UpdateFamilyParams{
		Name: util.RandomName(),
		ID:   family.ID.String(),
	}

	testCases := []struct {
		name          string
		params        UpdateFamilyParams
		stubs         func(store *databaseMock.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			params: params,
			stubs: func(store *databaseMock.MockStore) {
				store.EXPECT().
					UpdateFamily(mock.Anything, UpdateFamilyToDBUpdateFamily(params)).
					Times(1).Return(family, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "BadRequest",
			params: UpdateFamilyParams{},
			stubs: func(store *databaseMock.MockStore) {
				store.EXPECT().
					UpdateFamily(mock.Anything, mock.Anything).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "InternalServerError",
			params: params,
			stubs: func(store *databaseMock.MockStore) {
				store.EXPECT().
					UpdateFamily(mock.Anything, UpdateFamilyToDBUpdateFamily(params)).
					Times(1).Return(database.Family{}, pgx.ErrTxClosed)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			store := new(databaseMock.MockStore)
			server := newTestServer(t, store)

			tc.stubs(store)

			recorder := httptest.NewRecorder()
			url := "/families"

			data, err := encodeJSON(tc.params)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestDeleteFamily(t *testing.T) {
	family := randomFamily()

	testCases := []struct {
		name          string
		familyID      uuid.UUID
		stubs         func(store *databaseMock.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			familyID: family.ID,
			stubs: func(store *databaseMock.MockStore) {
				store.EXPECT().
					DeleteFamily(mock.Anything, family.ID).
					Times(1).Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:     "BadRequest",
			familyID: uuid.UUID{},
			stubs: func(store *databaseMock.MockStore) {
				store.EXPECT().
					DeleteFamily(mock.Anything, mock.Anything).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:     "NotFound",
			familyID: family.ID,
			stubs: func(store *databaseMock.MockStore) {
				store.EXPECT().
					DeleteFamily(mock.Anything, family.ID).
					Times(1).Return(pgx.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			store := new(databaseMock.MockStore)
			server := newTestServer(t, store)

			tc.stubs(store)

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/families/%s", tc.familyID.String())

			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
