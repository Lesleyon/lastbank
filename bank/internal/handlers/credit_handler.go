package handler

import (
	"database/sql"
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"bank/internal/middleware"
	"bank/internal/models"
	"bank/internal/repositories"
)

type CreditHandler struct {
	creditRepo  *repositories.CreditRepo
	paymentRepo *repositories.PaymentScheduleRepo
	accountRepo *repositories.AccountRepo
}

func NewCreditHandler(db *sql.DB) *CreditHandler {
	return &CreditHandler{
		creditRepo:  repositories.NewCreditRepo(db),
		paymentRepo: repositories.NewPaymentScheduleRepo(db),
		accountRepo: repositories.NewAccountRepo(db),
	}
}

func calculateAnnuityPayment(amount float64, annualRate float64, months int) float64 {
	monthlyRate := annualRate / 12 / 100
	if monthlyRate == 0 {
		return amount / float64(months)
	}
	annuity := monthlyRate * math.Pow(1+monthlyRate, float64(months)) / (math.Pow(1+monthlyRate, float64(months)) - 1)
	return amount * annuity
}

func (h *CreditHandler) CreateCredit(w http.ResponseWriter, r *http.Request) {
	userIDStr := middleware.GetUserID(r)
	userID, _ := strconv.ParseInt(userIDStr, 10, 64)

	var body struct {
		AccountID  int64   `json:"account_id"`
		Amount     float64 `json:"amount"`
		TermMonths int     `json:"term_months"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	rate := 15.0
	monthlyPayment := calculateAnnuityPayment(body.Amount, rate, body.TermMonths)

	credit := &models.Credit{
		UserID:         userID,
		AccountID:      body.AccountID,
		Amount:         body.Amount,
		InterestRate:   rate,
		MonthlyPayment: monthlyPayment,
		RemainingDebt:  body.Amount,
		IssuedAt:       time.Now(),
	}

	if err := h.creditRepo.Create(credit); err != nil {
		http.Error(w, "Failed to create credit", http.StatusInternalServerError)
		return
	}

	if err := h.accountRepo.UpdateBalance(body.AccountID, body.Amount); err != nil {
		http.Error(w, "Failed to update balance", http.StatusInternalServerError)
		return
	}

	monthlyRate := rate / 12 / 100
	remaining := body.Amount

	for i := 1; i <= body.TermMonths; i++ {
		interest := remaining * monthlyRate
		principal := monthlyPayment - interest
		if principal > remaining {
			principal = remaining
		}
		remaining -= principal

		schedule := &models.PaymentSchedule{
			CreditID:      credit.ID,
			PaymentNumber: i,
			DueDate:       time.Now().AddDate(0, i, 0),
			Amount:        monthlyPayment,
			Principal:     principal,
			Interest:      interest,
			Status:        "pending",
		}
		h.paymentRepo.Create(schedule)
	}

	json.NewEncoder(w).Encode(credit)
}

func (h *CreditHandler) GetSchedule(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	creditID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		http.Error(w, "Invalid credit ID", http.StatusBadRequest)
		return
	}

	schedules, err := h.paymentRepo.FindByCreditID(creditID)
	if err != nil {
		http.Error(w, "Failed to get schedule", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(schedules)
}
