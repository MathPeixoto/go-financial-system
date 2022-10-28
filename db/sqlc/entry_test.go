package db

import (
	"context"
	"database/sql"
	"github.com/MathPeixoto/go-financial-system/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandomEntry(t *testing.T, account Account) Entry {
	arg := CreateEntryParams{account.ID, util.RandomMoney()}
	entry, err := testQueries.CreateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestQueries_CreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	createRandomEntry(t, account)
}

func TestQueries_GetEntry(t *testing.T) {
	account := createRandomAccount(t)
	entryOne := createRandomEntry(t, account)

	entryTwo, err := testQueries.GetEntry(context.Background(), entryOne.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entryTwo)
	require.Equal(t, entryOne, entryTwo)
}

func TestQueries_UpdateEntry(t *testing.T) {
	account := createRandomAccount(t)
	entryOne := createRandomEntry(t, account)

	arg := UpdateEntryParams{
		ID:     entryOne.ID,
		Amount: util.RandomMoney(),
	}

	entryTwo, err := testQueries.UpdateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entryTwo)
	require.Equal(t, arg.Amount, entryTwo.Amount)
	require.Equal(t, entryOne.ID, entryTwo.ID)
	require.Equal(t, entryOne.AccountID, entryTwo.AccountID)
}

func TestQueries_DeleteEntry(t *testing.T) {
	account := createRandomAccount(t)
	entryOne := createRandomEntry(t, account)

	err := testQueries.DeleteEntry(context.Background(), entryOne.ID)

	require.NoError(t, err)

	entryTwo, err := testQueries.GetEntry(context.Background(), entryOne.ID)

	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, entryTwo)
}

func TestQueries_ListEntries(t *testing.T) {
	account := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomEntry(t, account)
	}

	arg := ListEntriesParams{
		Limit:  5,
		Offset: 5,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, entries, 5)
	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}

}
