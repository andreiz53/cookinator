package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	database "github.com/andreiz53/cookinator/database/handlers"
	databaseMock "github.com/andreiz53/cookinator/database/mocks"
	"github.com/andreiz53/cookinator/util"
)

func randomUser(t *testing.T) database.User {
	hashedPassword, err := util.HashPassword(util.RandomString(8))
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)
	return database.User{
		ID:        uuid.New(),
		FirstName: util.RandomFirstName(),
		Email:     util.RandomEmail(),
		Password:  hashedPassword,
	}

}

func TestCreateUser(t *testing.T) {
}

func TestGetUsers(t *testing.T) {

}

func TestGetUserByID(t *testing.T) {
	user := randomUser(t)

	store := new(databaseMock.MockStore)
	store.EXPECT().
		GetUserByID(mock.Anything, user.ID).
		Times(1).
		Return(user, nil)

	server := newTestServer(t, store)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/users/%s", user.ID.String())
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)
	server.router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusOK, recorder.Code)
}

func TestUpdateUserEmail(t *testing.T) {

}

func TestUpdateUserPassword(t *testing.T) {

}

func TestUpdateUserInfo(t *testing.T) {

}

func TestDeleteUser(t *testing.T) {

}
