package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	fmt.Println(">> before:", acc1.Balance, acc2.Balance)
	itr := 5
	amount := int64(10)
	errChan := make(chan error)
	resultChan := make(chan TransferTxResult)

	for i := 0; i < itr; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: acc1.ID,
				ToAccountID:   acc2.ID,
				Amount:        amount,
			})
			errChan <- err
			resultChan <- result
		}()
	}
	existed := make(map[int]bool)
	for i := 0; i < itr; i++ {
		err := <-errChan
		require.NoError(t, err)

		result := <-resultChan
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, acc1.ID, transfer.FromAccountID)
		require.Equal(t, acc2.ID, transfer.ToAccountID)
		require.Equal(t, transfer.Amount, amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		// check whether tranfer save in the transfer db or not
		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entry
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, acc1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		// check whether entry save in the entries db or not
		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		// check entry
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, acc2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		// check whether entry save in the entries db or not
		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, acc1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, acc2.ID, toAccount.ID)

		// check accounts balance
		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)
		diff1 := acc1.Balance - fromAccount.Balance
		diff2 := acc2.Balance - toAccount.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= itr)

		require.NotContains(t, existed, k)
		existed[k] = true

	}
	// check the final updated balances
	updateAccount1, err := testQueries.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)

	updateAccount2, err := testQueries.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)
	fmt.Println(">> after:", updateAccount1.Balance, updateAccount2.Balance)

	require.Equal(t, acc1.Balance-int64(itr)*amount, updateAccount1.Balance)
	require.Equal(t, acc2.Balance-int64(itr)*amount, updateAccount2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	fmt.Println(">> before:", acc1.Balance, acc2.Balance)
	itr := 10
	amount := int64(10)
	errChan := make(chan error)

	for i := 0; i < itr; i++ {
		fromAccountID := acc1.ID
		toAccountID := acc2.ID

		if i%2 == 0 {
			fromAccountID = acc2.ID
			toAccountID = acc1.ID
		}
		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errChan <- err
		}()
	}
	for i := 0; i < itr; i++ {
		err := <-errChan
		require.NoError(t, err)
	}
	// check the final updated balances
	updateAccount1, err := testQueries.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)

	updateAccount2, err := testQueries.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)
	fmt.Println(">> after:", updateAccount1.Balance, updateAccount2.Balance)

	require.Equal(t, acc1.Balance, updateAccount1.Balance)
	require.Equal(t, acc2.Balance, updateAccount2.Balance)
}
