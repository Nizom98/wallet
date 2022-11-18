package repository

import (
	"errors"
	"example.com/wallet/internal/models"
	"math/rand"
	"sync"
	"time"
)

var (
	errWalletNotFound = errors.New("wallet not found")
	seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// WalletRepository ...
type WalletRepository struct {
	// muWallets для конкурентного доступа к wallets
	muWallets *sync.RWMutex
	// wallets хранилище кошелков
	wallets []*wallet
}

// NewRepo конструктор репозитория
func NewRepo() *WalletRepository {
	return &WalletRepository{
		muWallets: new(sync.RWMutex),
		wallets: nil,
	}
}

// Transaction для конкурентной записи в хранилище.
// Сразу после вызова метода и до окончания доступ к хранилищу может иметь только один писатель.
func (repo *WalletRepository) Transaction(fn func(repo models.WalletRepository) error) error {
	repo.muWallets.Lock()
	defer repo.muWallets.Unlock()

	return fn(repo)
}

// Create создание кошелька.
func (repo *WalletRepository) Create(name string, balance float64, status bool) models.Walleter {
	newWallet := &wallet{
		id:      genNewID(),
		name:    name,
		balance: balance,
		status:  status,
	}

	repo.wallets = append(repo.wallets, newWallet)

	return newWallet
}

// ByID получаем кошелек по идентификатору.
// При отсутствии кошелка вернется ошибка errWalletNotFound.
func (repo *WalletRepository) ByID(id string) (models.Walleter, error) {
	for _, wal := range repo.wallets {
		if wal.id == id {
			return wal, nil
		}
	}

	return nil, errWalletNotFound
}


// All получение всего списка кошельков
func (repo *WalletRepository) All() []models.Walleter {
	walletList := make([]models.Walleter, 0, len(repo.wallets))
	for _, wallet := range repo.wallets {
		walletList = append(walletList, models.Walleter(wallet))
	}

	return walletList
}

// UpdateByID обновление данных кошелька.
// Все параметры(кроме id) являются опциональными.
// Если какой-то параметр отсутствует(равен nil), то данное поле не будет обновлено.
func (repo *WalletRepository) UpdateByID(id string, name *string, balance *float64, status *bool) error {
	pos := repo.walletPos(id)
	if pos == -1 {
		return errWalletNotFound
	}

	wal := repo.wallets[pos]

	if name != nil {
		wal.name = *name
	}
	if balance != nil {
		wal.balance = *balance
	}
	if status != nil {
		wal.status = *status
	}
	return nil
}

// walletPos определяем позицию(индекс) кошелька в хранилище
func (repo *WalletRepository) walletPos(id string) int {
	for pos, wal := range repo.wallets {
		if wal.id == id {
			return pos
		}
	}

	return -1
}

func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func genNewID() string {
	return stringWithCharset(8, charset)
}
