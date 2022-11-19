package api

import (
	db "github.com/MathPeixoto/go-financial-system/db/sqlc"
	"github.com/MathPeixoto/go-financial-system/util"
)

// Accounts

func getAccountParams(args CreateAccountRequest) db.CreateAccountParams {
	return db.CreateAccountParams{
		Currency: args.Currency,
		Balance:  0,
	}
}

func randomAccount(owner string) db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    owner,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func createAccount(accountParams db.CreateAccountParams) db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    accountParams.Owner,
		Balance:  accountParams.Balance,
		Currency: accountParams.Currency,
	}
}

func updatedAccount(account db.Account, accountBalanceParams db.AddAccountBalanceParams) db.Account {
	return db.Account{
		ID:       account.ID,
		Owner:    account.Owner,
		Balance:  account.Balance + accountBalanceParams.Amount,
		Currency: account.Currency,
	}
}

func brlAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    "BRL Account",
		Balance:  1000,
		Currency: "BRL",
	}
}

func usdAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    "BRL Account",
		Balance:  1000,
		Currency: "USD",
	}
}

// Transfers
func getTransferParams(request transferRequest) db.TransferTxParams {
	return db.TransferTxParams{
		FromAccountID: request.FromAccountID,
		ToAccountID:   request.ToAccountID,
		Amount:        request.Amount,
	}
}

func createTransferTx(transferTxParams db.TransferTxParams) db.TransferTxResult {
	return db.TransferTxResult{
		Transfer: db.Transfer{
			ID:            util.RandomInt(1, 1000),
			FromAccountID: transferTxParams.FromAccountID,
			ToAccountID:   transferTxParams.ToAccountID,
		},

		FromAccount: db.Account{
			ID:       transferTxParams.FromAccountID,
			Owner:    "BRL Account",
			Balance:  1000 - transferTxParams.Amount,
			Currency: "BRL",
		},
		ToAccount: db.Account{
			ID:       transferTxParams.ToAccountID,
			Owner:    "BRL Account",
			Balance:  1000 + transferTxParams.Amount,
			Currency: "BRL",
		},
		FromEntry: db.Entry{
			ID:        util.RandomInt(1, 1000),
			AccountID: transferTxParams.FromAccountID,
			Amount:    -transferTxParams.Amount,
		},
		ToEntry: db.Entry{
			ID:        util.RandomInt(1, 1000),
			AccountID: transferTxParams.ToAccountID,
			Amount:    transferTxParams.Amount,
		},
	}
}

func randomTransfer(accountOne, accountTwo db.Account) db.Transfer {
	return db.Transfer{
		ID:            util.RandomInt(1, 1000),
		FromAccountID: accountOne.ID,
		ToAccountID:   accountTwo.ID,
		Amount:        util.RandomMoney(),
	}
}
