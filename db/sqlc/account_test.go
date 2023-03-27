package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/ritik/simplebank/db/utils"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrencies(),
	}

	acc, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, arg.Balance, acc.Balance)
	require.Equal(t, arg.Currency, acc.Currency)
	require.Equal(t, arg.Owner, acc.Owner)

	return acc
}
func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	acc := createRandomAccount(t)
	account, err := testQueries.GetAccount(context.Background(), acc.ID)
	require.NoError(t, err)

	require.NotEmpty(t, account)
	require.Equal(t, acc.ID, account.ID)

	require.Equal(t, acc.Balance, account.Balance)
	require.Equal(t, acc.Currency, account.Currency)
	require.Equal(t, acc.Owner, account.Owner)
	require.WithinDuration(t, acc.CreatedAt, account.CreatedAt, time.Second)
}

func TestUpdate(t *testing.T) {
	acc := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      acc.ID,
		Balance: utils.RandomMoney(),
	}

	account, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, acc.ID, account.ID)
	require.Equal(t, account.Balance, arg.Balance)
	require.WithinDuration(t, acc.CreatedAt, account.CreatedAt, time.Second)
}

func TestDelete(t *testing.T) {
	acc := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), acc.ID)
	require.NoError(t, err)

	account, err := testQueries.GetAccount(context.Background(), acc.ID)
	require.Error(t, err)
	require.Empty(t, account)
	require.Equal(t, err, sql.ErrNoRows.Error())
}

func TestListAccount(t *testing.T) {

	for i := 0; i < 6; i++ {
		createRandomAccount(t)
	}
	arg := ListAccountsParams{
		Limit:  2,
		Offset: 3,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 2)

	for _, r := range accounts {
		require.NotEmpty(t, r)
	}
}
