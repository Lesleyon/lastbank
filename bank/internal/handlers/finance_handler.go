package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"bank/internal/middleware"
	"bank/internal/repositories"
)

type FinancesHandler struct {
	accountRepo *repositories.AccountRepo
	txRepo      *repositories.TransactionRepo
}

func NewFinancesHandler(db *sql.DB) *FinancesHandler {
	return &FinancesHandler{
		accountRepo: repositories.NewAccountRepo(db),
		txRepo:      repositories.NewTransactionRepo(db),
	}
}

func (h *FinancesHandler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	userIDStr := middleware.GetUserID(r)
	userID, _ := strconv.ParseInt(userIDStr, 10, 64)

	accounts, err := h.accountRepo.FindByUserID(userID)
	if err != nil {
		http.Error(w, "Failed to get analytics", http.StatusInternalServerError)
		return
	}

	var totalBalance float64
	for _, acc := range accounts {
		totalBalance += acc.Balance
	}

	stats := map[string]interface{}{
		"total_balance":  totalBalance,
		"accounts_count": len(accounts),
	}

	json.NewEncoder(w).Encode(stats)
}

func (h *FinancesHandler) PredictBalance(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	accountID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	days, err := strconv.Atoi(parts[4])
	if err != nil {
		http.Error(w, "Invalid days", http.StatusBadRequest)
		return
	}

	if days > 365 {
		days = 365
	}

	account, err := h.accountRepo.FindByID(accountID)
	if err != nil {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	predictions := make([]float64, days)
	for i := 0; i < days; i++ {
		predictions[i] = account.Balance - float64(i)*100
		if predictions[i] < 0 {
			predictions[i] = 0
		}
	}

	json.NewEncoder(w).Encode(predictions)
}
