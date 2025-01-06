package httpv1

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sumup-oss/go-pkgs/logger"

	"encoding/json"
	"errors"
	"tribe-payments-wallet-golang-interview-assignment/internal/wallet"

	"github.com/go-chi/chi/v5"
)

type CreateWalletRequest struct {
	Currency string `json:"currency"`
}

type WalletResponse struct {
	ID        string  `json:"id"`
	Balance   float64 `json:"balance"`
	Currency  string  `json:"currency"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type OperationRequest struct {
	Balance float64 `json:"balance"`
}

func NewCreateWalletHandler(svc wallet.Service, log logger.StructuredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			WriteError(w, http.StatusBadRequest, "Request body is required")
			return
		}

		var req CreateWalletRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("Failed to decode request body")
			WriteError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if req.Currency == "" {
			WriteError(w, http.StatusBadRequest, "Currency is required")
			return
		}

		req.Currency = strings.ToUpper(req.Currency)
		if len(req.Currency) != 3 {
			WriteError(w, http.StatusBadRequest, "Currency must be a 3-letter ISO code")
			return
		}

		newWallet, err := svc.CreateWallet(r.Context(), req.Currency)
		if err != nil {
			log.Error("Failed to create wallet")
			WriteError(w, http.StatusInternalServerError, "Failed to create wallet")
			return
		}

		response := WalletResponse{
			ID:        newWallet.ID,
			Balance:   newWallet.Balance,
			Currency:  newWallet.Currency,
			CreatedAt: newWallet.CreatedAt.Format(time.RFC3339),
			UpdatedAt: newWallet.UpdatedAt.Format(time.RFC3339),
		}

		WriteJSON(w, http.StatusCreated, response)
	}
}

func NewGetWalletHandler(svc wallet.Service, log logger.StructuredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		walletID := chi.URLParam(r, "id")
		if walletID == "" {
			WriteError(w, http.StatusBadRequest, "Wallet ID is required")
			return
		}

		foundWallet, err := svc.GetWallet(r.Context(), walletID)
		if err != nil {
			// Add this debug line
			log.Error(fmt.Sprintf("Error getting wallet: %v", err))

			if errors.Is(err, wallet.ErrWalletNotFound) {
				WriteError(w, http.StatusNotFound, "Wallet not found")
				return
			}
			log.Error("Failed to get wallet")
			WriteError(w, http.StatusInternalServerError, "Failed to get wallet")
			return
		}

		response := WalletResponse{
			ID:        foundWallet.ID,
			Balance:   foundWallet.Balance,
			Currency:  foundWallet.Currency,
			CreatedAt: foundWallet.CreatedAt.Format(time.RFC3339),
			UpdatedAt: foundWallet.UpdatedAt.Format(time.RFC3339),
		}

		WriteJSON(w, http.StatusOK, response)
	}
}

func NewDepositHandler(svc wallet.Service, log logger.StructuredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		walletID := chi.URLParam(r, "id")
		if walletID == "" {
			WriteError(w, http.StatusBadRequest, "Wallet ID is required")
			return
		}

		var req OperationRequest
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields() // Prevent unknown fields

		if err := decoder.Decode(&req); err != nil {
			log.Error(fmt.Sprintf("Failed to decode deposit request body: %v", err))
			WriteError(w, http.StatusBadRequest, "Invalid JSON format or unknown fields")
			return
		}

		if req.Balance <= 0 {
			WriteError(w, http.StatusBadRequest, "Amount must be greater than zero")
			return
		}

		if err := svc.Deposit(r.Context(), walletID, req.Balance); err != nil {
			switch {
			case errors.Is(err, wallet.ErrInvalidAmount):
				WriteError(w, http.StatusBadRequest, "Invalid amount")
			case errors.Is(err, wallet.ErrWalletNotFound):
				WriteError(w, http.StatusNotFound, "Wallet not found")
			default:
				log.Error("Failed to process deposit")
				WriteError(w, http.StatusInternalServerError, "Failed to process deposit")
			}
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func NewWithdrawHandler(svc wallet.Service, log logger.StructuredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		walletID := chi.URLParam(r, "id")
		if walletID == "" {
			WriteError(w, http.StatusBadRequest, "Wallet ID is required")
			return
		}

		var req OperationRequest
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields() // Prevent unknown fields

		if err := decoder.Decode(&req); err != nil {
			log.Error(fmt.Sprintf("Failed to decode withdraw request body: %v", err))
			WriteError(w, http.StatusBadRequest, "Invalid JSON format or unknown fields")
			return
		}

		if req.Balance <= 0 {
			WriteError(w, http.StatusBadRequest, "Amount must be greater than zero")
			return
		}

		if err := svc.Withdraw(r.Context(), walletID, req.Balance); err != nil {
			switch {
			case errors.Is(err, wallet.ErrInvalidAmount):
				WriteError(w, http.StatusBadRequest, "Invalid amount")
			case errors.Is(err, wallet.ErrInsufficientFunds):
				WriteError(w, http.StatusBadRequest, "Insufficient funds")
			case errors.Is(err, wallet.ErrWalletNotFound):
				WriteError(w, http.StatusNotFound, "Wallet not found")
			default:
				log.Error("Failed to process withdrawal")
				WriteError(w, http.StatusInternalServerError, "Failed to process withdrawal")
			}
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, map[string]string{"error": message})
}
