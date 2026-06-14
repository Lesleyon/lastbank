package repositories

import (
	"bank/internal/models"
	"database/sql"
)

type TransactionRepo struct {
	db *sql.DB
}

func NewTransactionRepo(db *sql.DB) *TransactionRepo {
	return &TransactionRepo{db: db}
}

func (r *TransactionRepo) Create(tx *models.Transaction) error {
	query := `INSERT INTO transactions (from_account_id, to_account_id, amount, type, status) 
              VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return r.db.QueryRow(query, tx.FromAccountID, tx.ToAccountID, tx.Amount, tx.Type, tx.Status).Scan(&tx.ID)
}
