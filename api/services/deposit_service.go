package services

import (
	"context"
	"database/sql"
	"errors"

	db "github.com/hhow09/simple_bank/db/sqlc"
	"github.com/hhow09/simple_bank/util"
)

type DepositService struct {
	store db.Store
}

func NewDepositService(store db.Store) DepositService {
	return DepositService{
		store: store,
	}
}

type DepositParams struct {
	User     db.User
	Currency string
	Amount   int64
}

func (dc *DepositService) Deposit(ctx context.Context, params DepositParams) (db.Transfer, error) {
	var transfer db.Transfer
	if !util.IsSupportedCurrency(params.Currency) {
		return transfer, errors.New("currency not supportted")
	}
	if params.Amount <= 0 {
		return transfer, errors.New("deposit should be more than 0")
	}
	user, err := dc.store.GetUser(ctx, params.User.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return transfer, errors.New("user not found")
		}
		return transfer, errors.New("internal server error")
	}
	return dc.store.CreateDepositTx(ctx, db.CreateDepositTxParams{User: user, Amount: params.Amount, Currency: params.Currency})
}
