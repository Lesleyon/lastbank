package models

import "time"

type PaymentSchedule struct {
	ID            int64      `json:"id"`
	CreditID      int64      `json:"credit_id"`
	PaymentNumber int        `json:"payment_number"`
	DueDate       time.Time  `json:"due_date"`
	Amount        float64    `json:"amount"`
	Principal     float64    `json:"principal"`
	Interest      float64    `json:"interest"`
	Status        string     `json:"status"`
	PaidAt        *time.Time `json:"paid_at"`
}
