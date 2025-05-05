package db

import (
	"context"
	"testing"

	"github.com/hhow09/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func TestCreateAccountTx(t *testing.T) {
	store := NewStore(testDB)
	user := createRandomUser(t)
	args := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
		AccType:  AccountTypeBank,
	}
	acc, err := store.CreateAccountTx(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, acc)

	accExt, err := store.GetExtAccount(context.Background(), GetExtAccountParams{Owner: user.Username, Currency: acc.Currency})
	require.NoError(t, err)
	require.NotEmpty(t, accExt)
	require.Equal(t, accExt.Owner, acc.Owner)
	require.Equal(t, accExt.Balance, int64(0))
	require.Equal(t, accExt.Currency, acc.Currency)
	require.Equal(t, accExt.AccType, AccountTypeExternal)
}
