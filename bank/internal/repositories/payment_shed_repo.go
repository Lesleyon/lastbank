package repositories

import (
	"bank/internal/models"
	"database/sql"
)

type PaymentScheduleRepo struct {
	db *sql.DB
}

func NewPaymentScheduleRepo(db *sql.DB) *PaymentScheduleRepo {
	return &PaymentScheduleRepo{db: db}
}

func (r *PaymentScheduleRepo) Create(s *models.PaymentSchedule) error {
	query := `INSERT INTO payment_schedules (credit_id, payment_number, due_date, amount, principal, interest, status) 
              VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	return r.db.QueryRow(query, s.CreditID, s.PaymentNumber, s.DueDate, s.Amount, s.Principal, s.Interest, s.Status).Scan(&s.ID)
}

func (r *PaymentScheduleRepo) FindByCreditID(creditID int64) ([]models.PaymentSchedule, error) {
	rows, err := r.db.Query(`SELECT id, credit_id, payment_number, due_date, amount, principal, interest, status, paid_at 
                              FROM payment_schedules WHERE credit_id = $1 ORDER BY payment_number`, creditID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.PaymentSchedule
	for rows.Next() {
		var s models.PaymentSchedule
		rows.Scan(&s.ID, &s.CreditID, &s.PaymentNumber, &s.DueDate, &s.Amount, &s.Principal, &s.Interest, &s.Status, &s.PaidAt)
		schedules = append(schedules, s)
	}
	return schedules, nil
}
