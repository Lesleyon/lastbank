package models

type Card struct {
	ID          int64  `json:"id"`
	AccountID   int64  `json:"account_id"`
	CardNumber  string `json:"card_number"`
	ExpiryMonth int    `json:"expiry_month"`
	ExpiryYear  int    `json:"expiry_year"`
	CVVHash     string `json:"-"`
	HMAC        string `json:"-"`
}
