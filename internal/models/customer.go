package models

import (
	"time"
)

type Customer struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (m *DBModel) InsertCustomer(c Customer) (int, error) {
	var id int
	query := `
		INSERT INTO customers (first_name, last_name, email)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err := m.DB.QueryRow(query, c.FirstName, c.LastName, c.Email).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
