package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/bflannery/banking/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	newAccountArgs := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoneyAmount(),
		Currency: util.RandomCurrency(),
	}

	newAccount, err := testQueries.CreateAccount(context.Background(), newAccountArgs)
	require.NoError(t, err)
	require.NotEmpty(t, newAccount)

	require.Equal(t, newAccount.Owner, newAccountArgs.Owner)
	require.Equal(t, newAccount.Balance, newAccountArgs.Balance)
	require.Equal(t, newAccount.Currency, newAccountArgs.Currency)

	require.NotZero(t, newAccount.ID)
	require.NotZero(t, newAccount.CreatedAt)

	return newAccount
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	newAccount := createRandomAccount(t)
	accountRecord, err := testQueries.GetAccount(context.Background(), newAccount.ID)
	require.NoError(t, err)
	require.NotEmpty(t, accountRecord)

	require.Equal(t, newAccount.ID, accountRecord.ID)
	require.Equal(t, newAccount.Owner, accountRecord.Owner)
	require.Equal(t, newAccount.Balance, accountRecord.Balance)
	require.Equal(t, newAccount.Currency, accountRecord.Currency)
	require.WithinDuration(t, newAccount.CreatedAt, accountRecord.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	newAccount := createRandomAccount(t)

	updateAccountArgs := UpdateAccountParams{
		ID:      newAccount.ID,
		Balance: util.RandomMoneyAmount(),
	}

	updatedAccount, err := testQueries.UpdateAccount(context.Background(), updateAccountArgs)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount)

	require.Equal(t, updatedAccount.ID, updateAccountArgs.ID)
	require.Equal(t, updatedAccount.Owner, newAccount.Owner)
	require.Equal(t, updatedAccount.Balance, updateAccountArgs.Balance)
	require.Equal(t, updatedAccount.Currency, newAccount.Currency)
	require.WithinDuration(t, updatedAccount.CreatedAt, newAccount.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	newAccount := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), newAccount.ID)
	require.NoError(t, err)

	accountRecord, err := testQueries.GetAccount(context.Background(), newAccount.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, accountRecord)
}

func TestListAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	listAccountsParams := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accountRecords, err := testQueries.ListAccounts(context.Background(), listAccountsParams)
	require.NoError(t, err)
	require.Len(t, accountRecords, 5)

	for _, accountRecord := range accountRecords {
		require.NotEmpty(t, accountRecord)
	}
}
