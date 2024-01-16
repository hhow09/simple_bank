package db

import "context"

func (store *SQLStore) CreateAccountTx(ctx context.Context, param CreateAccountParams) (Account, error) {
	// 1. create account
	// 2. create external account (double entry)
	var acc Account
	err := store.execTx(ctx, func(q *Queries) error {
		var err_ error
		acc, err_ = q.CreateAccount(ctx, param)
		if err_ != nil {
			return err_
		}
		param_ext := CreateAccountParams{
			Owner:    param.Owner,
			Currency: param.Currency,
			Balance:  0,
			AccType:  AccountTypeExternal,
		}
		_, err_ = q.CreateAccount(ctx, param_ext)
		if err_ != nil {
			return err_
		}
		return nil
	})

	return acc, err
}
