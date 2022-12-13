package util

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestPassword(t *testing.T) {
	password := RandomString(6)

	hashOne, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashOne)

	err = CheckPasswordHash(password, hashOne)
	require.NoError(t, err)

	wrongPassword := RandomString(6)
	err = CheckPasswordHash(wrongPassword, hashOne)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hashTwo, err := HashPassword(password)

	require.NoError(t, err)
	require.NotEmpty(t, hashTwo)
	require.NotEqualf(t, hashOne, hashTwo, "hashes should not be equal")
}
