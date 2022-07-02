package db

import (
	"context"
	"testing"
	"time"

	"github.com/bflannery/banking/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {

	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	newUserArgs := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	newUser, err := testQueries.CreateUser(context.Background(), newUserArgs)
	require.NoError(t, err)
	require.NotEmpty(t, newUser)

	require.Equal(t, newUser.Username, newUserArgs.Username)
	require.Equal(t, newUser.HashedPassword, newUserArgs.HashedPassword)
	require.Equal(t, newUser.FullName, newUserArgs.FullName)
	require.Equal(t, newUser.Email, newUserArgs.Email)

	require.True(t, newUser.PasswordChangedAt.IsZero())
	require.NotZero(t, newUser.CreatedAt)

	return newUser
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	newUser := createRandomUser(t)
	userRecord, err := testQueries.GetUser(context.Background(), newUser.Username)
	require.NoError(t, err)
	require.NotEmpty(t, userRecord)

	require.Equal(t, newUser.Username, userRecord.Username)
	require.Equal(t, newUser.HashedPassword, userRecord.HashedPassword)
	require.Equal(t, newUser.FullName, userRecord.FullName)
	require.Equal(t, newUser.Email, userRecord.Email)
	require.WithinDuration(t, newUser.PasswordChangedAt, userRecord.PasswordChangedAt, time.Second)
	require.WithinDuration(t, newUser.CreatedAt, userRecord.CreatedAt, time.Second)
}
