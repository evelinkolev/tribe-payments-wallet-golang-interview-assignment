package api

import (
	"tribe-payments-wallet-golang-interview-assignment/internal/api/httpv1"
	"tribe-payments-wallet-golang-interview-assignment/internal/wallet"

	"github.com/go-chi/chi/v5"
	"github.com/sumup-oss/go-pkgs/logger"
)

func RegisterRoutes(
	mux *chi.Mux,
	log logger.StructuredLogger,
	walletService wallet.Service,
) {
	mux.Get("/live", Health)

	mux.Route("/v1", func(r chi.Router) {
		r.Post("/wallets", httpv1.NewCreateWalletHandler(walletService, log))
		r.Get("/wallets/{id}", httpv1.NewGetWalletHandler(walletService, log))
		r.Post("/wallets/{id}/deposit", httpv1.NewDepositHandler(walletService, log))
		r.Post("/wallets/{id}/withdraw", httpv1.NewWithdrawHandler(walletService, log))
	})
}
