package db

import "context"

type CreateDepositTxParams struct {
	User     User
	Currency string
	Amount   int64 `json:"amount"`
}

func (store *SQLStore) CreateDepositTx(ctx context.Context, param CreateDepositTxParams) (Transfer, error) {
	var res TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err_ error

		acc, err_ := q.GetAccountByCurrencyForUpdate(ctx, GetAccountByCurrencyForUpdateParams{Owner: param.User.Username, Currency: param.Currency})
		if err_ != nil {
			return err_
		}
		extacc, err_ := q.GetExtAccountForUpdate(ctx, GetExtAccountForUpdateParams{Owner: param.User.Username, Currency: param.Currency})
		if err_ != nil {
			return err_
		}
		res, err_ = TransferWithTx(q, ctx, TransferTxParams{FromAccountID: extacc.ID, ToAccountID: acc.ID, Amount: param.Amount})
		if err_ != nil {
			return err_
		}
		return nil
	})

	return res.Transfer, err
}
