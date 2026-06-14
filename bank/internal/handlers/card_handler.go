package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"bank/internal/models"
	"bank/internal/repositories"

	"golang.org/x/crypto/bcrypt"
)

type CardHandler struct {
	cardRepo    *repositories.CardRepo
	accountRepo *repositories.AccountRepo
	hmacSecret  []byte
}

func NewCardHandler(db *sql.DB, hmacSecret string) *CardHandler {
	return &CardHandler{
		cardRepo:    repositories.NewCardRepo(db),
		accountRepo: repositories.NewAccountRepo(db),
		hmacSecret:  []byte(hmacSecret),
	}
}

func generateCardNumber() string {
	return "4532" + fmt.Sprintf("%012d", time.Now().UnixNano()%1000000000000)
}

func (h *CardHandler) CreateCard(w http.ResponseWriter, r *http.Request) {
	accountIDStr := r.URL.Query().Get("account_id")
	if accountIDStr == "" {
		http.Error(w, "account_id required", http.StatusBadRequest)
		return
	}
	accountID, _ := strconv.ParseInt(accountIDStr, 10, 64)

	cardNumber := generateCardNumber()
	expiryMonth := 12
	expiryYear := 2028
	cvv := fmt.Sprintf("%03d", time.Now().Nanosecond()%1000)

	cvvHash, _ := bcrypt.GenerateFromPassword([]byte(cvv), bcrypt.DefaultCost)

	hmacHash := hmac.New(sha256.New, h.hmacSecret)
	hmacHash.Write([]byte(cardNumber))
	hmacValue := hex.EncodeToString(hmacHash.Sum(nil))

	card := &models.Card{
		AccountID:   accountID,
		CardNumber:  cardNumber,
		ExpiryMonth: expiryMonth,
		ExpiryYear:  expiryYear,
		CVVHash:     string(cvvHash),
		HMAC:        hmacValue,
	}

	if err := h.cardRepo.Create(card); err != nil {
		http.Error(w, "Failed to create card", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(card)
}

func (h *CardHandler) GetCards(w http.ResponseWriter, r *http.Request) {
	accountIDStr := r.URL.Query().Get("account_id")
	if accountIDStr == "" {
		http.Error(w, "account_id required", http.StatusBadRequest)
		return
	}
	accountID, _ := strconv.ParseInt(accountIDStr, 10, 64)

	cards, err := h.cardRepo.FindByAccountID(accountID)
	if err != nil {
		http.Error(w, "Failed to get cards", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(cards)
}
