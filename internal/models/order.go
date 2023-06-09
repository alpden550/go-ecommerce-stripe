package models

import (
	"context"
	"time"
)

type OrderStatus struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type Order struct {
	ID             int          `json:"id"`
	WidgetID       int          `json:"widget_id"`
	SubscriptionID int          `json:"subscription_id"`
	TransactionID  int          `json:"transaction_id"`
	CustomerID     int          `json:"customer_id"`
	StatusID       int          `json:"status_id"`
	Quantity       int          `json:"quantity"`
	Amount         int          `json:"amount"`
	CreatedAt      time.Time    `json:"-"`
	UpdatedAt      time.Time    `json:"-"`
	Widget         Widget       `json:"widget"`
	Subscription   Subscription `json:"subscription"`
	Transaction    Transaction  `json:"transaction"`
	Customer       Customer     `json:"customer"`
	Status         OrderStatus  `json:"status"`
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

func (m *DBModel) InsertSubscriptionOrder(order Order) (int, error) {
	var id int

	query := `
		INSERT INTO orders (subscription_id, transaction_id, customer_id, status_id, quantity, amount)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	err := m.DB.QueryRow(
		query, order.SubscriptionID, order.TransactionID, order.CustomerID, order.StatusID, order.Quantity, order.Amount,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *DBModel) GetWidgetOrders() ([]*Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var orders []*Order
	query := `
		SELECT
    		o.id, o.quantity, o.amount, o.created_at, o.status_id,
    		w.name, w.description,
    		t.id, t.currency, t.last_four, t.expire_year, t.expire_month, t.payment_intent_code,
    		s.name, c.email, c.first_name, c.last_name
		FROM orders o
		LEFT JOIN widgets w ON w.id = o.widget_id
		LEFT JOIN transactions t ON t.id = o.transaction_id
		LEFT JOIN statuses s on s.id = o.status_id
		LEFT JOIN customers c on c.id = o.customer_id
		WHERE o.widget_id IS NOT NULL
		ORDER BY o.created_at DESC
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var o Order
		if err = rows.Scan(
			&o.ID,
			&o.Quantity,
			&o.Amount,
			&o.CreatedAt,
			&o.StatusID,
			&o.Widget.Name,
			&o.Widget.Description,
			&o.Transaction.ID,
			&o.Transaction.Currency,
			&o.Transaction.LastFour,
			&o.Transaction.ExpireYear,
			&o.Transaction.ExpireMonth,
			&o.Transaction.PaymentIntentCode,
			&o.Status.Name,
			&o.Customer.Email,
			&o.Customer.FirstName,
			&o.Customer.LastName,
		); err != nil {
			return orders, err
		}

		orders = append(orders, &o)
	}

	if err = rows.Err(); err != nil {
		return orders, err
	}

	return orders, nil
}

func (m *DBModel) GetSubscriptionsOrders() ([]*Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var orders []*Order
	query := `
		SELECT
    		o.id, o.quantity, o.amount, o.created_at,
    		sc.name, sc.description,
    		t.id, t.currency, t.last_four, t.expire_year, t.expire_month,
    		s.name, c.email, c.first_name, c.last_name
		FROM orders o
		LEFT JOIN subscriptions sc on sc.id = o.subscription_id
		LEFT JOIN transactions t ON t.id = o.transaction_id
		LEFT JOIN statuses s on s.id = o.status_id
		LEFT JOIN customers c on c.id = o.customer_id
		WHERE o.subscription_id IS NOT NULL
		ORDER BY o.created_at DESC
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var o Order
		if err = rows.Scan(
			&o.ID,
			&o.Quantity,
			&o.Amount,
			&o.CreatedAt,
			&o.Subscription.Name,
			&o.Subscription.Description,
			&o.Transaction.ID,
			&o.Transaction.Currency,
			&o.Transaction.LastFour,
			&o.Transaction.ExpireYear,
			&o.Transaction.ExpireMonth,
			&o.Status.Name,
			&o.Customer.Email,
			&o.Customer.FirstName,
			&o.Customer.LastName,
		); err != nil {
			return orders, err
		}

		orders = append(orders, &o)
	}

	if err = rows.Err(); err != nil {
		return orders, err
	}

	return orders, nil
}

func (m *DBModel) GetWidgetOrderByID(id int) (*Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var o Order

	query := `
		SELECT
    		o.id, o.quantity, o.amount, o.created_at, o.status_id,
    		w.name, w.description,
    		t.id, t.currency, t.last_four, t.expire_year, t.expire_month, t.payment_intent_code,
    		s.name, c.email, c.first_name, c.last_name
		FROM orders o
		LEFT JOIN widgets w on w.id = o.widget_id
		LEFT JOIN transactions t ON t.id = o.transaction_id
		LEFT JOIN statuses s on s.id = o.status_id
		LEFT JOIN customers c on c.id = o.customer_id
		WHERE o.id=$1
	`

	if err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&o.ID,
		&o.Quantity,
		&o.Amount,
		&o.CreatedAt,
		&o.StatusID,
		&o.Widget.Name,
		&o.Widget.Description,
		&o.Transaction.ID,
		&o.Transaction.Currency,
		&o.Transaction.LastFour,
		&o.Transaction.ExpireYear,
		&o.Transaction.ExpireMonth,
		&o.Transaction.PaymentIntentCode,
		&o.Status.Name,
		&o.Customer.Email,
		&o.Customer.FirstName,
		&o.Customer.LastName,
	); err != nil {
		return nil, err
	}

	return &o, nil
}

func (m *DBModel) GetSubscriptionOrderByID(id int) (*Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var o Order

	query := `
		SELECT
    		o.id, o.quantity, o.amount, o.created_at,
    		sc.name, sc.description,
    		t.id, t.currency, t.last_four, t.expire_year, t.expire_month,
    		s.name, c.email, c.first_name, c.last_name
		FROM orders o
		LEFT JOIN subscriptions sc on sc.id = o.subscription_id
		LEFT JOIN transactions t ON t.id = o.transaction_id
		LEFT JOIN statuses s on s.id = o.status_id
		LEFT JOIN customers c on c.id = o.customer_id
		WHERE o.id=$1
	`

	if err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&o.ID,
		&o.Quantity,
		&o.Amount,
		&o.CreatedAt,
		&o.Subscription.Name,
		&o.Subscription.Description,
		&o.Transaction.ID,
		&o.Transaction.Currency,
		&o.Transaction.LastFour,
		&o.Transaction.ExpireYear,
		&o.Transaction.ExpireMonth,
		&o.Status.Name,
		&o.Customer.Email,
		&o.Customer.FirstName,
		&o.Customer.LastName,
	); err != nil {
		return nil, err
	}

	return &o, nil
}
