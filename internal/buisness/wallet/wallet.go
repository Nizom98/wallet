package wallet

import (
	"errors"
	"fmt"

	"github.com/Nizom98/wallet/internal/models"
	"github.com/Nizom98/wallet/internal/utils"
)

const (
	defaultBalance = 0
	defaultStatus  = true
)

var (
	errAmountLessThanOne = errors.New("amount cannot be less than 1")
	errEmptyName         = errors.New("empty wallet name")
	errNotEnoughBalance  = errors.New("wallet has not enough balance")
	errSameWallet        = errors.New("same wallet")
)

type manager struct {
	repo models.WalletRepository
}

// NewManager конструктор менеджера кошельков
func NewManager(repo models.WalletRepository) *manager {
	return &manager{
		repo: repo,
	}
}

// Create создаем новый кошелек.
// name - наименование кошелька(не пустое).
func (man *manager) Create(name string) (models.Walleter, error) {
	if name == "" {
		return nil, errEmptyName
	}
	var newWallet models.Walleter
	err := man.repo.Transaction(func(repo models.WalletRepository) error {
		newWallet = repo.Create(name, defaultBalance, defaultStatus)
		return nil
	})

	return newWallet, err
}

// ByID получаем кошелек по идентификатору(даже если деактивирован).
func (man *manager) ByID(id string) (models.Walleter, error) {
	return man.repo.ByID(id)
}

// List получаем весь список кошельков(включая деактивированные)
func (man *manager) List() []models.Walleter {
	return man.repo.All()
}

// IncreaseBalanceBy пополнение кошелька.
// id - какой кошелек пополняем.
// amount - сумма пополнения(больше 0).
func (man *manager) IncreaseBalanceBy(id string, amount float64) error {
	if amount <= 0 {
		return errAmountLessThanOne
	}

	errTx := man.repo.Transaction(func(repo models.WalletRepository) error {
		wallet, err := repo.ByID(id)
		if err != nil {
			return fmt.Errorf("wallet %s: %w", id, err)
		}

		newBalance := wallet.Balance() + amount
		return repo.UpdateByID(id, nil, utils.Ptr[float64](newBalance), nil)
	})

	return errTx
}

// DecreaseBalanceBy снятие средств из кошелька.
// id - из какого кошелька снимаем.
// amount - сумма снятия(больше 0).
func (man *manager) DecreaseBalanceBy(id string, amount float64) error {
	if amount <= 0 {
		return errAmountLessThanOne
	}

	errTx := man.repo.Transaction(func(repo models.WalletRepository) error {
		wallet, err := repo.ByID(id)
		if err != nil {
			return fmt.Errorf("wallet %s: %w", id, err)
		}

		newBalance := wallet.Balance() - amount
		if newBalance < 0 {
			return fmt.Errorf("wallet %s: %w", wallet.ID(), errNotEnoughBalance)
		}

		return repo.UpdateByID(id, nil, utils.Ptr[float64](newBalance), nil)
	})

	return errTx
}

// TransferBalance перевод средств из одного кошелька в другой.
// Перевод в рамках одного кошелька запрещена.
// fromID - из какого кошелька переводи.
// toID - в какой кошелек переводим.
// amount - сумма перевода(больше 0).
func (man *manager) TransferBalance(fromID, toID string, amount float64) error {
	if fromID == toID {
		return errSameWallet
	}
	if amount <= 0 {
		return errAmountLessThanOne
	}

	errTx := man.repo.Transaction(func(repo models.WalletRepository) error {
		fromWallet, err := repo.ByID(fromID)
		if err != nil {
			return fmt.Errorf("cannot get source wallet by id %s: %w", fromID, err)
		}
		toWallet, err := repo.ByID(toID)
		if err != nil {
			return fmt.Errorf("cannot get dest wallet by id %s: %w", toID, err)
		}

		if fromWallet.Balance() < amount {
			return fmt.Errorf("wallet %s: %w", fromWallet.ID(), errNotEnoughBalance)
		}

		err = repo.UpdateByID(fromID, nil, utils.Ptr[float64](fromWallet.Balance()-amount), nil)
		if err != nil {
			return fmt.Errorf("cannot update source wallet: %w", err)
		}

		err = repo.UpdateByID(toID, nil, utils.Ptr[float64](toWallet.Balance()+amount), nil)
		if err != nil {
			return fmt.Errorf("cannot update dest wallet: %w", err)
		}
		return nil
	})

	return errTx
}

// DeactivateByID деактивируем кошелек по идентификатору.
func (man *manager) DeactivateByID(id string) error {
	errTx := man.repo.Transaction(func(repo models.WalletRepository) error {
		wallet, err := repo.ByID(id)
		if err != nil {
			return fmt.Errorf("cannot get wallet by id %s: %w", id, err)
		}

		err = repo.UpdateByID(wallet.ID(), nil, nil, utils.Ptr[bool](false))
		if err != nil {
			return fmt.Errorf("cannot update dest wallet: %w", err)
		}
		return nil
	})

	return errTx
}

// UpdateName обновляем наименование кошелька.
// Пустое наименование не допускается.
func (man *manager) UpdateName(id, name string) error {
	if name == "" {
		return errEmptyName
	}
	errTx := man.repo.Transaction(func(repo models.WalletRepository) error {
		err := repo.UpdateByID(id, utils.Ptr[string](name), nil, nil)
		if err != nil {
			return fmt.Errorf("cannot update dest wallet: %w", err)
		}
		return nil
	})

	return errTx
}
