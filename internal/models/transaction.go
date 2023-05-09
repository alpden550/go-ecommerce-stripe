package models

import (
	"time"
)

type TransactionStatus struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type Transaction struct {
	ID                  int       `json:"id"`
	Amount              int       `json:"amount"`
	Currency            string    `json:"currency"`
	LastFour            string    `json:"last_four"`
	ExpireMonth         int       `json:"expire_month"`
	ExpireYear          int       `json:"expire_year"`
	BankReturnCode      string    `json:"bank_return_code"`
	PaymentMethodCode   string    `json:"payment_method_code"`
	PaymentIntentCode   string    `json:"payment_intent_code"`
	TransactionStatusID int       `json:"transaction_status_id"`
	CreatedAt           time.Time `json:"-"`
	UpdatedAt           time.Time `json:"-"`
}

func (m *DBModel) InsertTransaction(t Transaction) (int, error) {
	var id int

	query := `
		INSERT INTO transactions 
		    (amount,
		     currency,
		     last_four,
		     bank_return_code,
		     transaction_status_id,
		     expire_month,
		     expire_year,
		     payment_method_code,
		     payment_intent_code
		     )
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`

	err := m.DB.QueryRow(
		query,
		t.Amount,
		t.Currency,
		t.LastFour,
		t.BankReturnCode,
		t.TransactionStatusID,
		t.ExpireMonth,
		t.ExpireYear,
		t.PaymentMethodCode,
		t.PaymentIntentCode,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
