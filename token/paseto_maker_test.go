package token

import (
	"testing"
	"time"

	"github.com/hhow09/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func newTestNewPasetoMaker() (Maker, error) {
	return NewPasetoMaker(util.Config{TokenSymmetricKey: util.RandomString(32)})
}

func TestNewPasetoMaker(t *testing.T) {
	maker, err := newTestNewPasetoMaker()
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	issueAt := time.Now()
	expireAt := issueAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issueAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expireAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := newTestNewPasetoMaker()
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	token, err := maker.CreateToken(username, -duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpireToken.Error())
	require.Nil(t, payload)
}

func TestInvalidToken(t *testing.T) {
	maker, err := newTestNewPasetoMaker()
	require.NoError(t, err)

	payload, err := maker.VerifyToken(util.RandomString(10))
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
