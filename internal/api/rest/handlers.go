package rest

import (
	"encoding/json"
	"example.com/wallet/internal/models"
	"github.com/gorilla/mux"
	"net/http"
)

func NewHandler(manWallet models.WalletManager, repoWallet models.WalletRepository) (*Handler, error) {
	return &Handler{
		manWallet: manWallet,
		repoWallet: repoWallet,
	}, nil
}

func (h *Handler) WalletCreateHandler(w http.ResponseWriter, req *http.Request) {
	dec := json.NewDecoder(req.Body)
	var request CreateWalletRequest
	err := dec.Decode(&request)
	if err != nil {
		printError(w, err.Error(), http.StatusBadRequest)
		return
	}
	wallet, err := h.manWallet.Create(request.Name)
	if err != nil {
		printError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	active := "active"
	if !wallet.Status() {
		active = "inactive"
	}

	resp := &CreateWalletResponse{
		ID:     wallet.ID(),
		Name:   wallet.Name(),
		Status: active,
	}
	printOk(w, resp)
}

func (h *Handler) WalletByIDHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	if id == "" {
		http.Error(w, "empty id", http.StatusBadRequest)
		return
	}

	wallet, err := h.manWallet.ByID(id)
	if err != nil {
		printError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := convertToWalletListResponse([]models.Walleter{wallet})
	printOk(w, resp[0])
}

func (h *Handler) WalletListHandler(w http.ResponseWriter, _ *http.Request) {
	wallets := h.manWallet.List()

	resp := convertToWalletListResponse(wallets)
	printOk(w, resp)
}

func (h *Handler) WalletDepositHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	if id == "" {
		http.Error(w, "empty id", http.StatusBadRequest)
		return
	}

	dec := json.NewDecoder(req.Body)
	var data WalletDepositWithdrawRequest
	err := dec.Decode(&data)
	if err != nil {
		printError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.manWallet.IncreaseBalanceBy(id, data.Amount)
	if err != nil {
		printError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	printOk(w, data)
}

func (h *Handler) WalletWithdrawHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	if id == "" {
		http.Error(w, "empty id", http.StatusBadRequest)
		return
	}

	dec := json.NewDecoder(req.Body)
	var data WalletDepositWithdrawRequest
	err := dec.Decode(&data)
	if err != nil {
		printError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.manWallet.DecreaseBalanceBy(id, data.Amount)
	if err != nil {
		printError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	printOk(w, data)
}

func (h *Handler) WalletTransferHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	if id == "" {
		http.Error(w, "empty id", http.StatusBadRequest)
		return
	}

	dec := json.NewDecoder(req.Body)
	var data WalletTransferRequest
	err := dec.Decode(&data)
	if err != nil {
		printError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.manWallet.TransferBalance(id, data.TransferTo, data.Amount)
	if err != nil {
		printError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	printOk(w, data)
}

func (h *Handler) WalletDeactivateHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	if id == "" {
		http.Error(w, "empty id", http.StatusBadRequest)
		return
	}

	err := h.manWallet.DeactivateByID(id)
	if err != nil {
		printError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	printOk(w, nil)
}

func (h *Handler) WalletUpdateHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	if id == "" {
		http.Error(w, "empty id", http.StatusBadRequest)
		return
	}

	dec := json.NewDecoder(req.Body)
	var data WalletUpdateNameRequest
	err := dec.Decode(&data)
	if err != nil {
		printError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.manWallet.UpdateName(id, data.Name)
	if err != nil {
		printError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	printOk(w, data)
}

func convertToWalletListResponse(inp []models.Walleter) []*WalletListResponse {
	out := make([]*WalletListResponse, 0, len(inp))

	for _, w := range inp {
		active := "active"
		if !w.Status() {
			active = "inactive"
		}
		out = append(out, &WalletListResponse{
			ID:      w.ID(),
			Name:    w.Name(),
			Balance: w.Balance(),
			Status:  active,
		})
	}

	return out
}

func printError(w http.ResponseWriter, err string, status int)  {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(
		&StatusResponse{
			Success: false,
			ErrMessage:    err,
		},
	)
}

func printOk(w http.ResponseWriter, data interface{})  {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		&StatusResponse{
			Success: true,
			Data: data,
		},
	)
}
