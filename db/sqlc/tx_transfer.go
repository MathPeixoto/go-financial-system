package db

import "context"

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx performs a money transfer from one account to the other
// It creates a transfer record, add account entries, and update accounts' balance within a single database transaction
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(queries *Queries) error {
		var err error

		result.Transfer, err = queries.CreateTransfer(ctx, CreateTransferParams(arg))
		if err != nil {
			return err
		}

		result.FromEntry, err = queries.CreateEntry(ctx, CreateEntryParams{
			arg.FromAccountID, -arg.Amount,
		})

		if err != nil {
			return err
		}

		result.ToEntry, err = queries.CreateEntry(ctx, CreateEntryParams{
			arg.ToAccountID, arg.Amount,
		})

		if err != nil {
			return err
		}

		// One good way to avoid deadlock is to update the account always in a given order.
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, _ = addMoney(
				ctx, queries, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount,
			)
		} else {
			result.ToAccount, result.FromAccount, _ = addMoney(
				ctx, queries, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount,
			)
		}

		return nil
	})

	return result, err
}

func addMoney(
	ctx context.Context, queries *Queries, accountIDOne int64, amountOne int64, accountIDTwo int64, amountTwo int64,
) (accountOne Account, accountTwo Account, err error) {
	accountOne, err = queries.AddAccountBalance(ctx, AddAccountBalanceParams{
		amountOne, accountIDOne,
	})
	if err != nil {
		return
	}

	accountTwo, err = queries.AddAccountBalance(ctx, AddAccountBalanceParams{
		amountTwo, accountIDTwo,
	})
	if err != nil {
		return
	}

	return
}
