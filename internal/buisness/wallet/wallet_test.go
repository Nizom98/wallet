package wallet

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Nizom98/wallet/internal/models"
	"github.com/stretchr/testify/assert"
)

//go:generate minimock -g -i github.com/Nizom98/wallet/internal/models.WalletRepository -o ./repository_mock_test.go -n RepositoryMock

func TestCreate(t *testing.T) {
	repo := NewRepositoryMock(t)
	man := NewManager(repo)
	expectName, expectID := "test_name", "test_id"

	repo.CreateMock.Set(func(name string, balance float64, status bool) (w1 models.Walleter) {
		return &fakeWallet{
			id:      expectID,
			name:    name,
			balance: balance,
			status:  status,
		}
	})
	repo.TransactionMock.Set(func(fn func(repo models.WalletRepository) error) (err error) {
		return fn(repo)
	})
	wallet, err := man.Create(expectName)

	assert.Nil(t, err)
	assert.True(t, expectID == wallet.ID())
	assert.True(t, expectName == wallet.Name())
	assert.True(t, defaultBalance == wallet.Balance())
	assert.True(t, defaultStatus == wallet.Status())
}

func TestIncreaseBalanceBy_found(t *testing.T) {
	repo := NewRepositoryMock(t)
	man := NewManager(repo)
	amount := float64(67)
	wallet := newFakeWallet("test_id", "test_name", defaultBalance)

	repo.ByIDMock.Return(wallet, nil)
	repo.UpdateByIDMock.Set(func(id string, name *string, balance *float64, status *bool) (err error) {
		assert.True(t, id == wallet.id)
		assert.Nil(t, name)
		assert.NotNil(t, balance)
		assert.True(t, *balance == defaultBalance+amount)
		assert.Nil(t, status)
		return nil
	})
	repo.TransactionMock.Set(func(fn func(repo models.WalletRepository) error) (err error) {
		return fn(repo)
	})

	err := man.IncreaseBalanceBy(wallet.id, amount)
	assert.Nil(t, err)
}

func TestIncreaseBalanceBy_notFound(t *testing.T) {
	repo := NewRepositoryMock(t)
	man := NewManager(repo)
	amount := float64(9999)
	expectErr := errors.New("not_found")
	unknownID := "test_id"

	repo.ByIDMock.Return(nil, expectErr)
	repo.TransactionMock.Set(func(fn func(repo models.WalletRepository) error) (err error) {
		return fn(repo)
	})

	err := man.IncreaseBalanceBy(unknownID, amount)
	assert.True(t, errors.Is(err, expectErr))
}

func TestIncreaseBalanceBy_incorrectAmount(t *testing.T) {
	man := NewManager(nil)
	incorrectAmount := float64(0)
	walletID := "test_id1"

	err := man.IncreaseBalanceBy(walletID, incorrectAmount)
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, errAmountLessThanOne))
}

func TestDecreaseBalanceBy_found(t *testing.T) {
	repo := NewRepositoryMock(t)
	man := NewManager(repo)
	amount := float64(67)
	oldBalance := float64(100)
	wallet := newFakeWallet("test_id", "test_name", oldBalance)

	repo.ByIDMock.Return(wallet, nil)
	repo.UpdateByIDMock.Set(func(id string, name *string, balance *float64, status *bool) (err error) {
		assert.True(t, id == wallet.id)
		assert.Nil(t, name)
		assert.NotNil(t, balance)
		assert.True(t, *balance == oldBalance-amount)
		assert.Nil(t, status)
		return nil
	})
	repo.TransactionMock.Set(func(fn func(repo models.WalletRepository) error) (err error) {
		return fn(repo)
	})

	err := man.DecreaseBalanceBy(wallet.id, amount)
	assert.Nil(t, err)
}

func TestDecreaseBalanceBy_notEnoughBalance(t *testing.T) {
	repo := NewRepositoryMock(t)
	man := NewManager(repo)
	amount := float64(9999)
	wallet := newFakeWallet("test_id", "test_name", 10)

	repo.ByIDMock.Return(wallet, nil)
	repo.TransactionMock.Set(func(fn func(repo models.WalletRepository) error) (err error) {
		return fn(repo)
	})

	err := man.DecreaseBalanceBy(wallet.id, amount)
	assert.True(t, errors.Is(err, errNotEnoughBalance))
}

func TestDecreaseBalanceBy_notFound(t *testing.T) {
	repo := NewRepositoryMock(t)
	man := NewManager(repo)
	amount := float64(9999)
	expectErr := errors.New("not_found")
	unknownID := "test_id"

	repo.ByIDMock.Return(nil, expectErr)
	repo.TransactionMock.Set(func(fn func(repo models.WalletRepository) error) (err error) {
		return fn(repo)
	})

	err := man.DecreaseBalanceBy(unknownID, amount)
	assert.True(t, errors.Is(err, expectErr))
}

func TestDecreaseBalanceBy_incorrectAmount(t *testing.T) {
	man := NewManager(nil)
	incorrectAmount := float64(0)
	walletID := "test_id1"

	err := man.DecreaseBalanceBy(walletID, incorrectAmount)
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, errAmountLessThanOne))
}

func TestTransferBalance_found(t *testing.T) {
	repo := NewRepositoryMock(t)
	man := NewManager(repo)
	amount := float64(100)

	fromWallet := newFakeWallet("from_id", "test_name_from", 999999)
	toWallet := newFakeWallet("to_id", "test_name_to", 0)

	repo.ByIDMock.Set(func(id string) (w1 models.Walleter, err error) {
		if id == fromWallet.id {
			return fromWallet, nil
		} else if id == toWallet.id {
			return toWallet, nil
		}
		return nil, fmt.Errorf("unexpected id")
	})
	repo.UpdateByIDMock.Set(func(id string, name *string, balance *float64, status *bool) (err error) {
		assert.Nil(t, name)
		assert.Nil(t, status)
		assert.NotNil(t, balance)

		if id == fromWallet.id {
			assert.True(t, *balance == fromWallet.balance-amount)
		} else if id == toWallet.id {
			assert.True(t, *balance == toWallet.balance+amount)
		} else {
			return fmt.Errorf("unexpected id")
		}

		return nil
	})
	repo.TransactionMock.Set(func(fn func(repo models.WalletRepository) error) (err error) {
		return fn(repo)
	})

	err := man.TransferBalance(fromWallet.id, toWallet.id, amount)
	assert.Nil(t, err)
}

func TestTransferBalance_notFound(t *testing.T) {
	repo := NewRepositoryMock(t)
	man := NewManager(repo)
	amount := float64(100)
	expectErr := errors.New("not_found")
	unknownID1 := "test_id_1"
	unknownID2 := "test_id_2"

	repo.ByIDMock.Return(nil, expectErr)
	repo.TransactionMock.Set(func(fn func(repo models.WalletRepository) error) (err error) {
		return fn(repo)
	})

	err := man.TransferBalance(unknownID1, unknownID2, amount)
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, expectErr))
}

func TestTransferBalance_sameWallet(t *testing.T) {
	man := NewManager(nil)
	amount := float64(100)
	walletID := "test_id"

	err := man.TransferBalance(walletID, walletID, amount)
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, errSameWallet))
}

func TestTransferBalance_incorrectAmount(t *testing.T) {
	man := NewManager(nil)
	incorrectAmount := float64(0)
	walletID1 := "test_id1"
	walletID2 := "test_id2"

	err := man.TransferBalance(walletID1, walletID2, incorrectAmount)
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, errAmountLessThanOne))
}

func TestDeactivateByID_found(t *testing.T) {
	repo := NewRepositoryMock(t)
	man := NewManager(repo)

	wallet := newFakeWallet("from_id", "test_name_from", 999999)

	repo.ByIDMock.Return(wallet, nil)
	repo.UpdateByIDMock.Set(func(id string, name *string, balance *float64, status *bool) (err error) {
		assert.Nil(t, name)
		assert.Nil(t, balance)
		assert.NotNil(t, status)
		assert.True(t, *status == false)
		return nil
	})
	repo.TransactionMock.Set(func(fn func(repo models.WalletRepository) error) (err error) {
		return fn(repo)
	})

	err := man.DeactivateByID(wallet.id)
	assert.Nil(t, err)
}

func TestDeactivateByID_notFound(t *testing.T) {
	repo := NewRepositoryMock(t)
	man := NewManager(repo)
	expectErr := errors.New("not_found")
	unknownID := "test_id_1"

	repo.ByIDMock.Return(nil, expectErr)
	repo.TransactionMock.Set(func(fn func(repo models.WalletRepository) error) (err error) {
		return fn(repo)
	})

	err := man.DeactivateByID(unknownID)
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, expectErr))
}

func newFakeWallet(id, name string, balance float64) *fakeWallet {
	return &fakeWallet{
		id:      id,
		name:    name,
		balance: balance,
		status:  defaultStatus,
	}
}

type fakeWallet struct {
	id      string  `json:"id"`
	name    string  `json:"name"`
	balance float64 `json:"balance"`
	status  bool    `json:"status"`
}

func (wal *fakeWallet) ID() string {
	return wal.id
}

func (wal *fakeWallet) Name() string {
	return wal.name
}

func (wal *fakeWallet) Balance() float64 {
	return wal.balance
}

func (wal *fakeWallet) Status() bool {
	return wal.status
}
