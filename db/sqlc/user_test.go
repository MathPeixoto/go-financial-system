package db

import (
	"context"
	"github.com/MathPeixoto/go-financial-system/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		util.RandomOwner(),
		"secret",
		util.RandomOwner(),
		util.RandomEmail(),
	}
	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.NotZero(t, user.CreatedAt)
	require.True(t, user.PasswordChangedAt.IsZero())

	return user
}

func TestQueries_CreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestQueries_GetUser(t *testing.T) {
	userOne := createRandomUser(t)

	userTwo, err := testQueries.GetUser(context.Background(), userOne.Username)

	require.NoError(t, err)
	require.NotEmpty(t, userTwo)
	require.Equal(t, userOne, userTwo)
}
