package notify

import (
	"example.com/wallet/internal/models"
	"log"
)

// NewManager конструктор уведомителя операций.
func NewManager(msgClient msgSender, manWallet models.WalletManager) *notify {
	return &notify{
		manWallet: manWallet,
		msgSender: msgClient,
	}
}

// Create перехватываем операцию создания и отправляем событие в брокер.
func (ntf *notify) Create(name string) (models.Walleter, error) {
	wallet, err := ntf.manWallet.Create(name)
	if err != nil {
		return wallet, err
	}

	ntf.sendEvent(eventWalletCreated)
	return wallet, nil
}

// ByID ...
func (ntf *notify) ByID(id string) (models.Walleter, error) {
	return ntf.manWallet.ByID(id)
}

// List ...
func (ntf *notify) List() []models.Walleter {
	return ntf.manWallet.List()
}

// IncreaseBalanceBy перехватываем операцию пополнения и отправляем событие в брокер.
func (ntf *notify) IncreaseBalanceBy(id string, amount float64) error {
	err := ntf.manWallet.IncreaseBalanceBy(id, amount)
	ntf.sendEvent(eventWalletDeposited)
	return err
}

// DecreaseBalanceBy перехватываем операцию снятия и отправляем событие в брокер.
func (ntf *notify) DecreaseBalanceBy(id string, amount float64) error {
	err := ntf.manWallet.DecreaseBalanceBy(id, amount)
	ntf.sendEvent(eventWalletWithdrawn)
	return err
}

// TransferBalance перехватываем операцию перевода и отправляем событие в брокер.
func (ntf *notify) TransferBalance(fromID, toID string, amount float64) error {
	err := ntf.manWallet.TransferBalance(fromID, toID, amount)
	ntf.sendEvent(eventWalletTransfered)
	return err
}

// DeactivateByID перехватываем операцию деактивации и отправляем событие в брокер.
func (ntf *notify) DeactivateByID(id string) error {
	err := ntf.manWallet.DeactivateByID(id)
	ntf.sendEvent(eventWalletDeleted)
	return err
}

// UpdateName ...
func (ntf *notify) UpdateName(id, name string) error {
	return ntf.manWallet.UpdateName(id, name)
}


// sendEvent отправляем сообщение брокеру.
// Если возникнет ошибка, то данные просто запишутся в лог.
func (ntf *notify) sendEvent(event string) {
	err := ntf.msgSender.Write([]byte(event))
	if err != nil {
		log.Printf("event %s not sent to nsq: %s\n", event, err.Error())
		return
	}

	log.Printf("event %s sent to nsq", event)
}