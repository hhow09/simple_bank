package db

import (
	"context"
	"testing"

	"github.com/hhow09/simple_bank/util"
	"github.com/stretchr/testify/require"
)

// create deposit:
// 1. create transfer
// 2. user's external account -= amount
// 3. user's bank account += amount
func TestCreateDepositTx(t *testing.T) {
	store := NewStore(testDB)
	user := createRandomUser(t)
	curr := util.RandomCurrency()
	args := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: curr,
		AccType:  AccountTypeBank,
	}
	acc, err := store.CreateAccountTx(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, acc)
	amt := util.RandomMoney()
	transfer, err := store.CreateDepositTx(context.Background(), CreateDepositTxParams{User: user, Currency: curr, Amount: amt})
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.Equal(t, transfer.Amount, amt)
	updatedAcc, err := store.GetAccount(context.Background(), acc.ID)
	require.NoError(t, err)
	extAcc, err := store.GetExtAccount(context.Background(), GetExtAccountParams{Owner: user.Username, Currency: curr})
	require.NoError(t, err)
	require.Equal(t, updatedAcc.Balance, acc.Balance+amt)
	require.Equal(t, extAcc.Balance, 0-amt)
}
