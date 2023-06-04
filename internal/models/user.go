package models

import (
	"context"
	"strings"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (m *DBModel) GetUserByEmail(email string) (User, error) {
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	email = strings.ToLower(email)
	query := `SELECT id, email, password, first_name, last_name  FROM users where email=$1`
	if err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
	); err != nil {
		return user, err
	}

	return user, nil
}

func (m *DBModel) SetUserPassword(user *User, hash string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `UPDATE users SET password = $1 WHERE id = $2`
	_, err := m.DB.ExecContext(ctx, query, hash, user.ID)
	if err != nil {
		return err
	}

	return nil
}
