package rest

import "example.com/wallet/internal/models"

type Handler struct {
	manWallet models.WalletManager
	manNotify models.WalletManager
	repoWallet models.WalletRepository
}

type CreateWalletRequest struct {
	Name string `json:"name"`
}

type CreateWalletResponse struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Status string `json:"status"`
}

type WalletListResponse struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Balance float64 `json:"balance"`
	Status string `json:"status"`
}

type StatusResponse struct {
	Success bool `json:"success"`
	ErrMessage string `json:"err_message,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

type WalletDepositWithdrawRequest struct {
	Amount float64 `json:"amount"`
}

type WalletTransferRequest struct {
	Amount float64 `json:"amount"`
	TransferTo string `json:"transfer_to"`
}

type WalletUpdateNameRequest struct {
	Name string `json:"name"`
}