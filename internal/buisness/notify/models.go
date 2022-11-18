package notify

import "example.com/wallet/internal/models"

const  (
	eventWalletCreated = "Wallet_Created"
	eventWalletDeleted = "Wallet_Deleted"
	eventWalletDeposited = "Wallet_Deposited"
	eventWalletWithdrawn = "Wallet_Withdrawn"
	eventWalletTransfered = "Wallet_Transfered"
)

type msgSender interface {
	Write(data []byte) error
}

type notify struct {
	manWallet models.WalletManager
	msgSender msgSender
}