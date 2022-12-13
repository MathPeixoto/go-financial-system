package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// SQLStore provides all functions to execute SQL queries and transactions
// It creates a composition so that I can extend the Queries functionalities without inheritance
type SQLStore struct {
	*Queries
	db *sql.DB
}

// NewStore create a new store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(queries *Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %w, rb error: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

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
