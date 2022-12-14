package db

import (
	"context"
	"database/sql"
	"github.com/MathPeixoto/go-financial-system/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)

	arg := CreateAccountParams{user.Username, util.RandomMoney(), util.RandomCurrency()}
	account, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestQueries_CreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestQueries_GetAccount(t *testing.T) {
	accountOne := createRandomAccount(t)

	accountTwo, err := testQueries.GetAccount(context.Background(), accountOne.ID)

	require.NoError(t, err)
	require.NotEmpty(t, accountTwo)
	require.Equal(t, accountOne, accountTwo)
}

func TestQueries_GetAccountByOwner(t *testing.T) {
	accountOne := createRandomAccount(t)

	accountTwo, err := testQueries.GetAccountByOwner(context.Background(), accountOne.Owner)

	require.NoError(t, err)
	require.NotEmpty(t, accountTwo)
	require.Equal(t, accountOne, accountTwo)
}

func TestQueries_GetAccountForUpdate(t *testing.T) {
	accountOne := createRandomAccount(t)

	accountTwo, err := testQueries.GetAccountForUpdate(context.Background(), accountOne.ID)

	require.NoError(t, err)
	require.NotEmpty(t, accountTwo)
	require.Equal(t, accountOne, accountTwo)
}

func TestQueries_UpdateAccount(t *testing.T) {
	accountOne := createRandomAccount(t)

	arg := UpdateAccountParams{accountOne.ID, util.RandomMoney()}
	accountTwo, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, accountTwo)
	require.Equal(t, arg.Balance, accountTwo.Balance)
	require.Equal(t, accountOne.ID, accountTwo.ID)
}

func TestQueries_DeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account1.ID)

	require.NoError(t, err)

	accountTwo, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, accountTwo)
}

func TestQueries_ListAccounts(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccount(t)
	}

	arg := ListAccountsParams{lastAccount.Owner, 5, 0}
	accounts, err := testQueries.ListAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.Owner, account.Owner)
	}
}
