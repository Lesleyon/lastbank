package models

import "time"

type Credit struct {
	ID             int64     `json:"id"`
	UserID         int64     `json:"user_id"`
	AccountID      int64     `json:"account_id"`
	Amount         float64   `json:"amount"`
	InterestRate   float64   `json:"interest_rate"`
	MonthlyPayment float64   `json:"monthly_payment"`
	RemainingDebt  float64   `json:"remaining_debt"`
	IssuedAt       time.Time `json:"issued_at"`
}
