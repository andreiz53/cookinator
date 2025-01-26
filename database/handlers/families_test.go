package database

import (
	"context"
	"testing"
	"time"

	"github.com/andreiz53/cookinator/util"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func createRandomFamily(t *testing.T) Family {
	user := createRandomUser(t)

	arg := CreateFamilyParams{
		Name:            util.RandomName(),
		CreatedByUserID: user.ID,
	}

	family, err := testQueries.CreateFamily(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, family)

	require.Equal(t, arg.Name, family.Name)
	require.Equal(t, arg.CreatedByUserID, family.CreatedByUserID)

	require.NotZero(t, family.ID)
	require.NotZero(t, family.CreatedAt)
	require.NotZero(t, family.UpdatedAt)

	return family
}

func TestCreateFamily(t *testing.T) {
	createRandomFamily(t)
}

func TestGetFamilyByID(t *testing.T) {
	family := createRandomFamily(t)

	family2, err := testQueries.GetFamilyByID(context.Background(), family.ID)
	require.NoError(t, err)
	require.NotEmpty(t, family2)

	require.Equal(t, family.ID, family2.ID)
	require.Equal(t, family.Name, family2.Name)
	require.Equal(t, family.CreatedByUserID, family2.CreatedByUserID)

	require.WithinDuration(t, family.CreatedAt.Time, family2.CreatedAt.Time, time.Second)
	require.WithinDuration(t, family.UpdatedAt.Time, family2.UpdatedAt.Time, time.Second)
}

func TestGetFamilyByUserID(t *testing.T) {
	family := createRandomFamily(t)

	family2, err := testQueries.GetFamilyByUserID(context.Background(), family.CreatedByUserID)
	require.NoError(t, err)
	require.NotEmpty(t, family2)

	require.Equal(t, family.ID, family2.ID)
	require.Equal(t, family.Name, family2.Name)
	require.Equal(t, family.CreatedByUserID, family2.CreatedByUserID)

	require.WithinDuration(t, family.CreatedAt.Time, family2.CreatedAt.Time, time.Second)
	require.WithinDuration(t, family.UpdatedAt.Time, family2.UpdatedAt.Time, time.Second)
}

func TestGetFamilies(t *testing.T) {
	createRandomFamily(t)
	createRandomFamily(t)
	createRandomFamily(t)

	families, err := testQueries.GetFamilies(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, families)

	require.True(t, len(families) >= 3)
}

func TestUpdateFamily(t *testing.T) {
	family := createRandomFamily(t)

	args := UpdateFamilyParams{
		ID:   family.ID,
		Name: util.RandomName(),
	}

	family2, err := testQueries.UpdateFamily(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, family2)

	require.Equal(t, family.ID, family2.ID)
	require.Equal(t, family.CreatedByUserID, family2.CreatedByUserID)
	require.Equal(t, args.Name, family2.Name)

	require.WithinDuration(t, family.CreatedAt.Time, family2.CreatedAt.Time, time.Second)
}

func TestDeleteFamily(t *testing.T) {
	family := createRandomFamily(t)

	err := testQueries.DeleteFamily(context.Background(), family.ID)
	require.NoError(t, err)

	family2, err := testQueries.GetFamilyByID(context.Background(), family.ID)
	require.Error(t, err)
	require.Empty(t, family2)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
}
