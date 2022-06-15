package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/bflannery/banking/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T) Transfer {
	// TODO: Randomize the to and from accounts

	newAccount1 := createRandomAccount(t)
	newAccount2 := createRandomAccount(t)

	newTransferArgs := CreateTransferParams{
		ToAccountID:   newAccount1.ID,
		FromAccountID: newAccount2.ID,
		Amount:        util.RandomMoneyAmount(),
	}

	newTransfer, err := testQueries.CreateTransfer(context.Background(), newTransferArgs)
	require.NoError(t, err)
	require.NotEmpty(t, newTransfer)

	require.Equal(t, newTransfer.ToAccountID, newTransferArgs.ToAccountID)
	require.Equal(t, newTransfer.Amount, newTransferArgs.Amount)

	require.NotZero(t, newTransfer.ID)
	require.NotZero(t, newTransfer.CreatedAt)

	return newTransfer
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	newTransfer := createRandomTransfer(t)
	transferRecord, err := testQueries.GetTransfer(context.Background(), newTransfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transferRecord)

	require.Equal(t, newTransfer.ID, transferRecord.ID)
	require.Equal(t, newTransfer.Amount, transferRecord.Amount)
	require.WithinDuration(t, newTransfer.CreatedAt, transferRecord.CreatedAt, time.Second)
}

func TestUpdateTransfer(t *testing.T) {
	newTransfer := createRandomTransfer(t)

	updateTransferArgs := UpdateTransferParams{
		ID:            newTransfer.ID,
		ToAccountID:   newTransfer.FromAccountID,
		FromAccountID: newTransfer.ToAccountID,
		Amount:        util.RandomMoneyAmount(),
	}

	updatedTransfer, err := testQueries.UpdateTransfer(context.Background(), updateTransferArgs)
	require.NoError(t, err)
	require.NotEmpty(t, updatedTransfer)

	require.Equal(t, updatedTransfer.ID, updateTransferArgs.ID)
	require.Equal(t, updatedTransfer.ToAccountID, updateTransferArgs.ToAccountID)
	require.Equal(t, updatedTransfer.FromAccountID, updateTransferArgs.FromAccountID)
	require.Equal(t, updatedTransfer.Amount, updateTransferArgs.Amount)
	require.WithinDuration(t, updatedTransfer.CreatedAt, newTransfer.CreatedAt, time.Second)
}

func TestDeleteTransfer(t *testing.T) {
	newTransfer := createRandomTransfer(t)
	err := testQueries.DeleteEntry(context.Background(), newTransfer.ID)
	require.NoError(t, err)

	transferRecord, err := testQueries.GetEntry(context.Background(), newTransfer.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, transferRecord)
}

func TestListTransfer(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomTransfer(t)
	}

	listTransfersParams := ListEntriesParams{
		Limit:  5,
		Offset: 5,
	}

	transferRecords, err := testQueries.ListEntries(context.Background(), listTransfersParams)
	require.NoError(t, err)
	require.Len(t, transferRecords, 5)

	for _, transferRecord := range transferRecords {
		require.NotEmpty(t, transferRecord)
	}
}
