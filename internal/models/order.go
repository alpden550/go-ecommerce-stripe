package models

import (
	"time"
)

type OrderStatus struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type Order struct {
	ID            int       `json:"id"`
	WidgetID      int       `json:"widget_id"`
	TransactionID int       `json:"transaction_id"`
	CustomerID    int       `json:"customer_id"`
	StatusID      int       `json:"status_id"`
	Quantity      int       `json:"quantity"`
	Amount        int       `json:"amount"`
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
}

func (m *DBModel) InsertOrder(order Order) (int, error) {
	var id int

	query := `
		INSERT INTO orders (widget_id, transaction_id, customer_id, status_id, quantity, amount)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	err := m.DB.QueryRow(
		query, order.WidgetID, order.TransactionID, order.CustomerID, order.StatusID, order.Quantity, order.Amount,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
