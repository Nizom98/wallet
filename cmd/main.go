package main

import (
	"example.com/wallet/internal/api/rest"
	"example.com/wallet/internal/buisness/notify"
	"example.com/wallet/internal/buisness/wallet"
	"example.com/wallet/internal/clients/nsq"
	"example.com/wallet/internal/repository"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	nsqTopic = "nsq_test"
	nsqTarget = "127.0.0.1:9999"
)


func main()  {
	nsq, err := nsq.NewClient(nsqTopic, nsqTarget)
	if err != nil {
		panic(err)
	}
	defer nsq.Stop()

	repoWallet := repository.NewRepo()
	manWallet := wallet.NewManager(repoWallet)
	manNotify := notify.NewManager(nsq, manWallet)

	handler, err := rest.NewHandler(manNotify, repoWallet)
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/wallet/", handler.WalletCreateHandler).Methods(http.MethodPost)
	r.HandleFunc("/wallets/{id}/", handler.WalletByIDHandler).Methods(http.MethodGet)
	r.HandleFunc("/wallets/", handler.WalletListHandler).Methods(http.MethodGet)
	r.HandleFunc("/wallets/{id}/", handler.WalletUpdateHandler).Methods(http.MethodPut)
	r.HandleFunc("/wallets/{id}/", handler.WalletDeactivateHandler).Methods(http.MethodDelete)
	r.HandleFunc("/wallets/{id}/deposit/", handler.WalletDepositHandler).Methods(http.MethodPost)
	r.HandleFunc("/wallets/{id}/withdraw/", handler.WalletWithdrawHandler).Methods(http.MethodPost)
	r.HandleFunc("/wallets/{id}/transfer/", handler.WalletTransferHandler).Methods(http.MethodPost)

	err = http.ListenAndServe(":80", r)
	panic(err)
}
