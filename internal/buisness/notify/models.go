package notify

import "github.com/Nizom98/wallet/internal/models"

const (
	eventWalletCreated    = "Wallet_Created"
	eventWalletDeleted    = "Wallet_Deleted"
	eventWalletDeposited  = "Wallet_Deposited"
	eventWalletWithdrawn  = "Wallet_Withdrawn"
	eventWalletTransfered = "Wallet_Transfered"
)

type msgSender interface {
	Write(data []byte) error
}

type notify struct {
	manWallet models.WalletManager
	msgSender msgSender
}

type eventData struct {
	Type   string  `json:"type"`
	Amount float64 `json:"amount"`
}
