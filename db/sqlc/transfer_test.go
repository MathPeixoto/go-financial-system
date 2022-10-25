package db

import (
	"bancario/util"
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandomTransfer(t *testing.T, accountOne, accountTwo Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: accountOne.ID,
		ToAccountID:   accountTwo.ID,
		Amount:        util.RandomMoney(),
	}

	Transfer, err := testQueries.CreateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, Transfer)
	require.Equal(t, arg.FromAccountID, Transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, Transfer.ToAccountID)
	require.Equal(t, arg.Amount, Transfer.Amount)
	require.NotZero(t, Transfer.ID)
	require.NotZero(t, Transfer.CreatedAt)

	return Transfer
}

func TestQueries_CreateTransfer(t *testing.T) {
	accountOne := createRandomAccount(t)
	accountTwo := createRandomAccount(t)
	createRandomTransfer(t, accountOne, accountTwo)
}

func TestQueries_GetTransfer(t *testing.T) {
	accountOne := createRandomAccount(t)
	accountTwo := createRandomAccount(t)
	transferOne := createRandomTransfer(t, accountOne, accountTwo)

	transferTwo, err := testQueries.GetTransfer(context.Background(), transferOne.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transferTwo)
	require.Equal(t, transferOne, transferTwo)
}

func TestQueries_UpdateTransfer(t *testing.T) {
	accountOne := createRandomAccount(t)
	accountTwo := createRandomAccount(t)
	transferOne := createRandomTransfer(t, accountOne, accountTwo)

	arg := UpdateTransferParams{
		ID:     transferOne.ID,
		Amount: util.RandomMoney(),
	}

	transferTwo, err := testQueries.UpdateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, transferTwo)
	require.Equal(t, arg.Amount, transferTwo.Amount)
	require.Equal(t, transferOne.ID, transferTwo.ID)
	require.Equal(t, transferOne.FromAccountID, transferTwo.FromAccountID)
	require.Equal(t, transferOne.ToAccountID, transferTwo.ToAccountID)
}

func TestQueries_DeleteTransfer(t *testing.T) {
	accountOne := createRandomAccount(t)
	accountTwo := createRandomAccount(t)
	transferOne := createRandomTransfer(t, accountOne, accountTwo)

	err := testQueries.DeleteTransfer(context.Background(), transferOne.ID)

	require.NoError(t, err)

	transferTwo, err := testQueries.GetTransfer(context.Background(), transferOne.ID)

	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, transferTwo)
}

func TestQueries_ListTransfers(t *testing.T) {
	accountOne := createRandomAccount(t)
	accountTwo := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomTransfer(t, accountOne, accountTwo)
	}

	arg := ListTransfersParams{
		Limit:  5,
		Offset: 5,
	}

	entries, err := testQueries.ListTransfers(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, entries, 5)
	for _, Transfer := range entries {
		require.NotEmpty(t, Transfer)
	}
}
