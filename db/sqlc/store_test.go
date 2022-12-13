package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStore_TransferTx(t *testing.T) {
	store := NewStore(testDB)
	accountOne := createRandomAccount(t)
	accountTwo := createRandomAccount(t)

	// run n concurrent result transactions
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				accountOne.ID, accountTwo.ID, amount,
			})

			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results

		// check result transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, accountOne.ID, transfer.FromAccountID)
		require.Equal(t, accountTwo.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, accountOne.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, accountTwo.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, accountOne.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, accountTwo.ID, toAccount.ID)

		// check accounts' balance
		diffOne := accountOne.Balance - fromAccount.Balance
		diffTwo := toAccount.Balance - accountTwo.Balance
		require.Equal(t, diffOne, diffTwo)
		require.True(t, diffOne > 0)
		require.True(t, diffTwo > 0)
		require.True(t, diffOne%amount == 0)

		k := int(diffOne / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	updateAccountOne, err := store.GetAccount(context.Background(), accountOne.ID)
	require.NoError(t, err)

	updateAccountTwo, err := store.GetAccount(context.Background(), accountTwo.ID)
	require.NoError(t, err)

	require.Equal(t, accountOne.Balance-int64(n)*amount, updateAccountOne.Balance)
	require.Equal(t, accountTwo.Balance+int64(n)*amount, updateAccountTwo.Balance)
}

func TestStore_TransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)
	accountOne := createRandomAccount(t)
	accountTwo := createRandomAccount(t)

	// run n concurrent result transactions
	n := 10
	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := accountOne.ID
		toAccountID := accountTwo.ID

		if i%2 == 1 {
			fromAccountID, toAccountID = accountTwo.ID, accountOne.ID
		}

		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				fromAccountID, toAccountID, amount,
			})

			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	updateAccountOne, err := store.GetAccount(context.Background(), accountOne.ID)
	require.NoError(t, err)

	updateAccountTwo, err := store.GetAccount(context.Background(), accountTwo.ID)
	require.NoError(t, err)

	require.Equal(t, accountOne.Balance, updateAccountOne.Balance)
	require.Equal(t, accountTwo.Balance, updateAccountTwo.Balance)
}
