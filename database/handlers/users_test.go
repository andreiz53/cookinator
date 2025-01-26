package database

import (
	"context"
	"testing"
	"time"

	"github.com/andreiz53/cookinator/util"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		FirstName: util.RandomFirstName(),
		Email:     util.RandomEmail(),
		Password:  util.RandomPassword(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.FirstName, user.FirstName)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Password, user.Password)

	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, user.UpdatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUserByID(t *testing.T) {
	user := createRandomUser(t)

	dbUser, err := testQueries.GetUserByID(context.Background(), user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, dbUser)

	require.Equal(t, user.ID, dbUser.ID)
	require.Equal(t, user.FirstName, dbUser.FirstName)
	require.Equal(t, user.Email, dbUser.Email)
	require.Equal(t, user.Password, dbUser.Password)
	require.Equal(t, user.FamilyID, dbUser.FamilyID)

	require.WithinDuration(t, user.CreatedAt.Time, dbUser.CreatedAt.Time, time.Second)
	require.WithinDuration(t, user.UpdatedAt.Time, dbUser.UpdatedAt.Time, time.Second)
}

func TestGetUserByEmail(t *testing.T) {
	user := createRandomUser(t)

	dbUser, err := testQueries.GetUserByEmail(context.Background(), user.Email)
	require.NoError(t, err)
	require.NotEmpty(t, dbUser)

	require.Equal(t, user.ID, dbUser.ID)
	require.Equal(t, user.FirstName, dbUser.FirstName)
	require.Equal(t, user.Email, dbUser.Email)
	require.Equal(t, user.Password, dbUser.Password)
	require.Equal(t, user.FamilyID, dbUser.FamilyID)

	require.WithinDuration(t, user.CreatedAt.Time, dbUser.CreatedAt.Time, time.Second)
	require.WithinDuration(t, user.UpdatedAt.Time, dbUser.UpdatedAt.Time, time.Second)
}

func TestGetUsers(t *testing.T) {
	createRandomUser(t)
	createRandomUser(t)
	createRandomUser(t)

	dbUsers, err := testQueries.GetUsers(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, dbUsers)

	require.True(t, len(dbUsers) >= 3)
}

func TestUpdateUserEmail(t *testing.T) {
	user := createRandomUser(t)

	arg := UpdateUserEmailParams{
		ID:    user.ID,
		Email: util.RandomEmail(),
	}

	newUser, err := testQueries.UpdateUserEmail(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, newUser)

	require.Equal(t, user.ID, newUser.ID)
	require.Equal(t, user.FirstName, newUser.FirstName)
	require.Equal(t, user.Password, newUser.Password)
	require.WithinDuration(t, user.CreatedAt.Time, newUser.CreatedAt.Time, time.Second)
	require.Equal(t, arg.Email, newUser.Email)
}

func TestUpdateUserInfo(t *testing.T) {
	user := createRandomUser(t)

	arg := UpdateUserInfoParams{
		ID:        user.ID,
		FirstName: util.RandomFirstName(),
	}

	newUser, err := testQueries.UpdateUserInfo(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, newUser)

	require.Equal(t, user.ID, newUser.ID)
	require.Equal(t, user.Password, newUser.Password)
	require.Equal(t, user.Email, newUser.Email)
	require.WithinDuration(t, user.CreatedAt.Time, newUser.CreatedAt.Time, time.Second)
	require.Equal(t, arg.FirstName, newUser.FirstName)
}

func TestUpdateUserPassword(t *testing.T) {
	user := createRandomUser(t)

	arg := UpdateUserPasswordParams{
		ID:       user.ID,
		Password: util.RandomPassword(),
	}

	newUser, err := testQueries.UpdateUserPassword(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, newUser)

	require.Equal(t, user.ID, newUser.ID)
	require.Equal(t, user.FirstName, newUser.FirstName)
	require.Equal(t, user.Email, newUser.Email)
	require.WithinDuration(t, user.CreatedAt.Time, newUser.CreatedAt.Time, time.Second)
	require.Equal(t, arg.Password, newUser.Password)
}

func TestDeleteUser(t *testing.T) {
	user := createRandomUser(t)

	err := testQueries.DeleteUser(context.Background(), user.ID)
	require.NoError(t, err)

	user2, err := testQueries.GetUserByID(context.Background(), user.ID)
	require.Error(t, err)
	require.Empty(t, user2)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
}
