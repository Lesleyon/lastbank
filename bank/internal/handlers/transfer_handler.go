package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"bank/internal/models"
	"bank/internal/repositories"
)

type TransferHandler struct {
	accountRepo *repositories.AccountRepo
	txRepo      *repositories.TransactionRepo
}

func NewTransferHandler(db *sql.DB) *TransferHandler {
	return &TransferHandler{
		accountRepo: repositories.NewAccountRepo(db),
		txRepo:      repositories.NewTransactionRepo(db),
	}
}

func (h *TransferHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	var body struct {
		FromAccountID int64   `json:"from_account_id"`
		ToAccountID   int64   `json:"to_account_id"`
		Amount        float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	from, err := h.accountRepo.FindByID(body.FromAccountID)
	if err != nil {
		http.Error(w, "From account not found", http.StatusNotFound)
		return
	}

	if from.Balance < body.Amount {
		http.Error(w, "Insufficient funds", http.StatusBadRequest)
		return
	}

	if err := h.accountRepo.UpdateBalance(body.FromAccountID, -body.Amount); err != nil {
		http.Error(w, "Transfer failed", http.StatusInternalServerError)
		return
	}

	if err := h.accountRepo.UpdateBalance(body.ToAccountID, body.Amount); err != nil {
		h.accountRepo.UpdateBalance(body.FromAccountID, body.Amount)
		http.Error(w, "Transfer failed", http.StatusInternalServerError)
		return
	}

	tx := &models.Transaction{
		FromAccountID: &body.FromAccountID,
		ToAccountID:   &body.ToAccountID,
		Amount:        body.Amount,
		Type:          "transfer",
		Status:        "completed",
	}
	h.txRepo.Create(tx)

	json.NewEncoder(w).Encode(map[string]string{"message": "Transfer successful"})
}
