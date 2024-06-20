package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hhow09/simple_bank/util"
	"go.uber.org/fx"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	CreateAccountTx(ctx context.Context, arg CreateAccountParams) (Account, error)
	CreateDepositTx(ctx context.Context, param CreateDepositTxParams) (Transfer, error)
}

// SQLStore provides all funcs to execute queries and transactions
type SQLStore struct {
	*Queries
	db *sql.DB
}

// openSQL connection
func openSQL(config util.Config) (*sql.DB, error) {
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		return nil, err
	}
	// ping to fail early
	err = conn.Ping()
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// creates a new Store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
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

func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	// 1. create transfer record
	// 2. create Entry of from account
	// 3. create Entry of to account
	// 4. update account
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result, err = TransferWithTx(q, ctx, arg)
		return err
	})

	return result, err
}

func TransferWithTx(q *Queries, ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var err error
	var result TransferTxResult
	result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
		FromAccountID: arg.FromAccountID,
		ToAccountID:   arg.ToAccountID,
		Amount:        arg.Amount,
	})
	if err != nil {
		return result, err
	}
	result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
		AccountID: arg.FromAccountID,
		Amount:    -arg.Amount,
	})
	if err != nil {
		return result, err
	}

	result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
		AccountID: arg.ToAccountID,
		Amount:    arg.Amount,
	})
	if err != nil {
		return result, err
	}

	// update accounts' balance
	if arg.FromAccountID < arg.ToAccountID {
		result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
	} else {
		result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
	}
	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	return
}

var Module = fx.Options(
	fx.Provide(openSQL),
	fx.Provide(NewStore),
)
