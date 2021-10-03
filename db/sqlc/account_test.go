package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/hhow09/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	args := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, args.Owner, account.Owner)
	require.Equal(t, args.Balance, account.Balance)
	require.Equal(t, args.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	acc := createRandomAccount(t)
	accGet, err := testQueries.GetAccount(context.Background(), acc.ID)
	require.NoError(t, err)
	require.NotEmpty(t, accGet)

	require.Equal(t, accGet.ID, acc.ID)
	require.Equal(t, accGet.Owner, acc.Owner)
	require.Equal(t, accGet.Balance, acc.Balance)
	require.Equal(t, accGet.Currency, acc.Currency)
	require.Equal(t, accGet.CreatedAt, acc.CreatedAt)
}

func TestUpdateAccount(t *testing.T) {
	acc := createRandomAccount(t)

	args := UpdateAccountParams{
		ID:      acc.ID,
		Balance: util.RandomMoney(),
	}
	accUpdated, err := testQueries.UpdateAccount(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, accUpdated)

	require.Equal(t, accUpdated.ID, acc.ID)
	require.Equal(t, accUpdated.Owner, acc.Owner)
	require.Equal(t, accUpdated.Balance, args.Balance)
	require.Equal(t, accUpdated.Currency, acc.Currency)
	require.Equal(t, accUpdated.CreatedAt, acc.CreatedAt)

}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	arg := ListAccountsParams{
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
