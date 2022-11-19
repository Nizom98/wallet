package notify

import (
	"encoding/json"

	"example.com/wallet/internal/models"
	log "github.com/sirupsen/logrus"
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

	ntf.sendEvent(eventWalletCreated, 0)
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
	ntf.sendEvent(eventWalletDeposited, amount)
	return err
}

// DecreaseBalanceBy перехватываем операцию снятия и отправляем событие в брокер.
func (ntf *notify) DecreaseBalanceBy(id string, amount float64) error {
	err := ntf.manWallet.DecreaseBalanceBy(id, amount)
	ntf.sendEvent(eventWalletWithdrawn, amount)
	return err
}

// TransferBalance перехватываем операцию перевода и отправляем событие в брокер.
func (ntf *notify) TransferBalance(fromID, toID string, amount float64) error {
	err := ntf.manWallet.TransferBalance(fromID, toID, amount)
	ntf.sendEvent(eventWalletTransfered, amount)
	return err
}

// DeactivateByID перехватываем операцию деактивации и отправляем событие в брокер.
func (ntf *notify) DeactivateByID(id string) error {
	err := ntf.manWallet.DeactivateByID(id)
	ntf.sendEvent(eventWalletDeleted, 0)
	return err
}

// UpdateName ...
func (ntf *notify) UpdateName(id, name string) error {
	return ntf.manWallet.UpdateName(id, name)
}

// sendEvent отправляем сообщение брокеру.
// Если возникнет ошибка, то данные запишутся в лог.
func (ntf *notify) sendEvent(eventType string, amount float64) {
	event := &eventData{
		Type:   eventType,
		Amount: amount,
	}
	bytes, err := json.Marshal(event)
	if err != nil {
		log.Errorf("err while marshaling event (type: %s, amount: %f): %s", event.Type, event.Amount, err.Error())
		return
	}
	err = ntf.msgSender.Write(bytes)
	if err != nil {
		log.Errorf("event (type: %s, amount: %f) NOT sent to nsq: %s", event.Type, event.Amount, err.Error())
		return
	}

	log.Errorf("event (type: %s, amount: %f) sent to nsq", event.Type, event.Amount)
}
