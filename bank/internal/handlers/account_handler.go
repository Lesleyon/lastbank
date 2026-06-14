package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"bank/internal/middleware"
	"bank/internal/models"
	"bank/internal/repositories"
)

type AccountHandler struct {
	accountRepo *repositories.AccountRepo
	txRepo      *repositories.TransactionRepo
}

func NewAccountHandler(db *sql.DB) *AccountHandler {
	return &AccountHandler{
		accountRepo: repositories.NewAccountRepo(db),
		txRepo:      repositories.NewTransactionRepo(db),
	}
}

func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	userIDStr := middleware.GetUserID(r)
	userID, _ := strconv.ParseInt(userIDStr, 10, 64)

	acc := &models.Account{UserID: userID}
	if err := h.accountRepo.Create(acc); err != nil {
		http.Error(w, "Failed to create account", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(acc)
}

func (h *AccountHandler) GetAccounts(w http.ResponseWriter, r *http.Request) {
	userIDStr := middleware.GetUserID(r)
	userID, _ := strconv.ParseInt(userIDStr, 10, 64)

	accounts, err := h.accountRepo.FindByUserID(userID)
	if err != nil {
		http.Error(w, "Failed to get accounts", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(accounts)
}

func (h *AccountHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	var body struct {
		AccountID int64   `json:"account_id"`
		Amount    float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.accountRepo.UpdateBalance(body.AccountID, body.Amount); err != nil {
		http.Error(w, "Deposit failed", http.StatusInternalServerError)
		return
	}

	tx := &models.Transaction{
		ToAccountID: &body.AccountID,
		Amount:      body.Amount,
		Type:        "deposit",
		Status:      "completed",
	}
	h.txRepo.Create(tx)

	json.NewEncoder(w).Encode(map[string]string{"message": "Deposit successful"})
}
