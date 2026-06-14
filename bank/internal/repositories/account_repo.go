package repositories

import (
	"bank/internal/models"
	"database/sql"
)

type AccountRepo struct {
	db *sql.DB
}

func NewAccountRepo(db *sql.DB) *AccountRepo {
	return &AccountRepo{db: db}
}

func (r *AccountRepo) Create(account *models.Account) error {
	query := `INSERT INTO accounts (user_id, balance) VALUES ($1, 0) RETURNING id`
	return r.db.QueryRow(query, account.UserID).Scan(&account.ID)
}

func (r *AccountRepo) FindByUserID(userID int64) ([]models.Account, error) {
	rows, err := r.db.Query("SELECT id, user_id, balance FROM accounts WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []models.Account
	for rows.Next() {
		var a models.Account
		rows.Scan(&a.ID, &a.UserID, &a.Balance)
		accounts = append(accounts, a)
	}
	return accounts, nil
}

func (r *AccountRepo) UpdateBalance(id int64, amount float64) error {
	_, err := r.db.Exec("UPDATE accounts SET balance = balance + $1 WHERE id = $2", amount, id)
	return err
}

func (r *AccountRepo) FindByID(id int64) (*models.Account, error) {
	a := &models.Account{}
	err := r.db.QueryRow("SELECT id, user_id, balance FROM accounts WHERE id = $1", id).Scan(&a.ID, &a.UserID, &a.Balance)
	return a, err
}
