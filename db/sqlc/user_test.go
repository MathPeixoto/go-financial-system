package db

import (
	"context"
	"database/sql"
	"github.com/MathPeixoto/go-financial-system/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		util.RandomOwner(),
		hashedPassword,
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

func TestQueries_UpdateFullNameUser(t *testing.T) {
	oldUser := createRandomUser(t)
	newFullName := util.RandomOwner()

	arg := UpdateUserParams{
		FullName: sql.NullString{String: newFullName, Valid: true},
		Username: oldUser.Username,
	}

	updatedUser, err := testQueries.UpdateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.Equal(t, arg.Username, updatedUser.Username)
	require.Equal(t, newFullName, updatedUser.FullName)
	require.NotZero(t, updatedUser.CreatedAt)
	require.True(t, updatedUser.PasswordChangedAt.IsZero())
}

func TestQueries_UpdateEmailUser(t *testing.T) {
	oldUser := createRandomUser(t)
	newEmail := util.RandomEmail()

	arg := UpdateUserParams{
		Email:    sql.NullString{String: newEmail, Valid: true},
		Username: oldUser.Username,
	}

	updatedUser, err := testQueries.UpdateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.Equal(t, arg.Username, updatedUser.Username)
	require.Equal(t, newEmail, updatedUser.Email)
	require.NotZero(t, updatedUser.CreatedAt)
	require.True(t, updatedUser.PasswordChangedAt.IsZero())
}

func TestQueries_UpdatePasswordUser(t *testing.T) {
	oldUser := createRandomUser(t)
	newPassword := util.RandomString(6)
	hashedPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)

	arg := UpdateUserParams{
		HashedPassword: sql.NullString{String: hashedPassword, Valid: true},
		Username:       oldUser.Username,
	}

	updatedUser, err := testQueries.UpdateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.Equal(t, arg.Username, updatedUser.Username)
	require.Equal(t, hashedPassword, updatedUser.HashedPassword)
	require.NotZero(t, updatedUser.CreatedAt)
	require.True(t, updatedUser.PasswordChangedAt.IsZero())
}
