package main

import (
	"net/http"

	"github.com/Nizom98/wallet/internal/api/rest"
	"github.com/Nizom98/wallet/internal/buisness/notify"
	"github.com/Nizom98/wallet/internal/buisness/wallet"
	"github.com/Nizom98/wallet/internal/clients/nsq"
	"github.com/Nizom98/wallet/internal/repository"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const (
	nsqTopic  = "nsq_test"
	nsqTarget = "127.0.0.1:9999"
	appAddr   = ":80"
	logLevel  = log.DebugLevel
)

func main() {
	log.SetLevel(logLevel)

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
	r.HandleFunc("/wallet/", handler.MiddlewareLog(handler.WalletCreateHandler)).Methods(http.MethodPost)
	r.HandleFunc("/wallets/{id}/", handler.MiddlewareLog(handler.WalletByIDHandler)).Methods(http.MethodGet)
	r.HandleFunc("/wallets/", handler.MiddlewareLog(handler.WalletListHandler)).Methods(http.MethodGet)
	r.HandleFunc("/wallets/{id}/", handler.MiddlewareLog(handler.WalletUpdateHandler)).Methods(http.MethodPut)
	r.HandleFunc("/wallets/{id}/", handler.MiddlewareLog(handler.WalletDeactivateHandler)).Methods(http.MethodDelete)
	r.HandleFunc("/wallets/{id}/deposit/", handler.MiddlewareLog(handler.WalletDepositHandler)).Methods(http.MethodPost)
	r.HandleFunc("/wallets/{id}/withdraw/", handler.MiddlewareLog(handler.WalletWithdrawHandler)).Methods(http.MethodPost)
	r.HandleFunc("/wallets/{id}/transfer/", handler.MiddlewareLog(handler.WalletTransferHandler)).Methods(http.MethodPost)

	log.Infof("app started on: %s", appAddr)
	defer func() {
		log.Infof("app finished on: %s", appAddr)
	}()
	err = http.ListenAndServe(appAddr, r)
	panic(err)
}
