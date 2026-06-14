package repositories

import (
	"bank/internal/models"
	"database/sql"
)

type CreditRepo struct {
	db *sql.DB
}

func NewCreditRepo(db *sql.DB) *CreditRepo {
	return &CreditRepo{db: db}
}

func (r *CreditRepo) Create(credit *models.Credit) error {
	query := `INSERT INTO credits (user_id, account_id, amount, interest_rate, monthly_payment, remaining_debt) 
              VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	return r.db.QueryRow(query, credit.UserID, credit.AccountID, credit.Amount,
		credit.InterestRate, credit.MonthlyPayment, credit.RemainingDebt).Scan(&credit.ID)
}

func (r *CreditRepo) FindByID(id int64) (*models.Credit, error) {
	c := &models.Credit{}
	query := `SELECT id, user_id, account_id, amount, interest_rate, monthly_payment, remaining_debt, issued_at 
              FROM credits WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&c.ID, &c.UserID, &c.AccountID, &c.Amount,
		&c.InterestRate, &c.MonthlyPayment, &c.RemainingDebt, &c.IssuedAt)
	return c, err
}
