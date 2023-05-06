package models

import (
	"context"
	"time"
)

type Widget struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Image          string    `json:"image"`
	InventoryLevel int       `json:"inventory_level"`
	Price          int       `json:"price"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}

func (m *DBModel) GetWidget(id int) (Widget, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var widget Widget
	query := `SELECT id, name, description, price, inventory_level, coalesce(image, ''), created_at FROM widgets WHERE id=$1`

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&widget.ID,
		&widget.Name,
		&widget.Description,
		&widget.Price,
		&widget.InventoryLevel,
		&widget.Image,
		&widget.CreatedAt,
	)
	if err != nil {
		return widget, err
	}

	return widget, nil
}