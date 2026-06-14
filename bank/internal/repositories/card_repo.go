package repositories

import (
	"bank/internal/models"
	"database/sql"
)

type CardRepo struct {
	db *sql.DB
}

func NewCardRepo(db *sql.DB) *CardRepo {
	return &CardRepo{db: db}
}

func (r *CardRepo) Create(card *models.Card) error {
	query := `INSERT INTO cards (account_id, card_number, expiry_month, expiry_year, cvv_hash, hmac) 
              VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	return r.db.QueryRow(query, card.AccountID, card.CardNumber, card.ExpiryMonth,
		card.ExpiryYear, card.CVVHash, card.HMAC).Scan(&card.ID)
}

func (r *CardRepo) FindByAccountID(accountID int64) ([]models.Card, error) {
	rows, err := r.db.Query("SELECT id, account_id, card_number, expiry_month, expiry_year FROM cards WHERE account_id = $1", accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []models.Card
	for rows.Next() {
		var c models.Card
		rows.Scan(&c.ID, &c.AccountID, &c.CardNumber, &c.ExpiryMonth, &c.ExpiryYear)
		cards = append(cards, c)
	}
	return cards, nil
}
