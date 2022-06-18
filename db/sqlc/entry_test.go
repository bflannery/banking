package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/bflannery/banking/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) Entry {
	newAccount := createRandomAccount(t)
	newEntryArgs := CreateEntryParams{
		AccountID: newAccount.ID,
		Amount:    util.RandomMoneyAmount(),
	}

	newEntry, err := testQueries.CreateEntry(context.Background(), newEntryArgs)
	require.NoError(t, err)
	require.NotEmpty(t, newEntry)

	require.Equal(t, newEntry.AccountID, newEntryArgs.AccountID)
	require.Equal(t, newEntry.Amount, newEntryArgs.Amount)

	require.NotZero(t, newAccount.ID)
	require.NotZero(t, newAccount.CreatedAt)

	return newEntry
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	newEntry := createRandomEntry(t)
	entryRecord, err := testQueries.GetEntry(context.Background(), newEntry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entryRecord)

	require.Equal(t, newEntry.ID, entryRecord.ID)
	require.Equal(t, newEntry.Amount, entryRecord.Amount)
	require.WithinDuration(t, newEntry.CreatedAt, entryRecord.CreatedAt, time.Second)
}

func TestUpdateEntry(t *testing.T) {
	newEntry := createRandomEntry(t)

	updateEntryArgs := UpdateEntryParams{
		ID:     newEntry.ID,
		Amount: util.RandomMoneyAmount(),
	}

	updatedEntry, err := testQueries.UpdateEntry(context.Background(), updateEntryArgs)
	require.NoError(t, err)
	require.NotEmpty(t, updatedEntry)

	require.Equal(t, updatedEntry.ID, updateEntryArgs.ID)
	require.Equal(t, updatedEntry.Amount, updateEntryArgs.Amount)
	require.WithinDuration(t, updatedEntry.CreatedAt, updatedEntry.CreatedAt, time.Second)
}

func TestDeleteEntry(t *testing.T) {
	newEntry := createRandomEntry(t)
	err := testQueries.DeleteEntry(context.Background(), newEntry.ID)
	require.NoError(t, err)

	entryRecord, err := testQueries.GetEntry(context.Background(), newEntry.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, entryRecord)
}

func TestListEntry(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomEntry(t)
	}

	listEntriesParams := ListEntriesParams{
		Limit:  5,
		Offset: 5,
	}

	entryRecords, err := testQueries.ListEntries(context.Background(), listEntriesParams)
	require.NoError(t, err)
	require.Len(t, entryRecords, 5)

	for _, entryRecord := range entryRecords {
		require.NotEmpty(t, entryRecord)
	}
}
