package models

import (
	"context"
	"time"
)

type Subscription struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	InventoryLevel int       `json:"inventory_level"`
	Price          int       `json:"price"`
	IsRecurring    bool      `json:"is_recurring"`
	PlanID         string    `json:"plan_id"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}

func (m *DBModel) GetSubscriptionByName(name string) (Subscription, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var subscription Subscription
	query := `
			SELECT 
    			id, name, description, price, inventory_level, created_at, is_recurring, plan_id 
			FROM subscriptions WHERE name=$1
			`

	row := m.DB.QueryRowContext(ctx, query, name)
	err := row.Scan(
		&subscription.ID,
		&subscription.Name,
		&subscription.Description,
		&subscription.Price,
		&subscription.InventoryLevel,
		&subscription.CreatedAt,
		&subscription.IsRecurring,
		&subscription.PlanID,
	)
	if err != nil {
		return subscription, err
	}

	return subscription, nil
}
