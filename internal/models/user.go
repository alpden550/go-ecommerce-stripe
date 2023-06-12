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
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
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

func (m *DBModel) GetAllUsers() ([]*User, error) {
	var users []*User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, email, first_name, last_name, created_at FROM users ORDER BY created_at DESC`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user User
		if err = rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.CreatedAt,
		); err != nil {
			return users, err
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return users, err
	}

	return users, nil
}

func (m *DBModel) GetUserByID(id int) (*User, error) {
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, email, first_name, last_name, created_at  FROM users where id=$1`
	if err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *DBModel) AddUser(user *User, hash string) (int, error) {
	var id int
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		INSERT INTO
			users (first_name, last_name, email, password)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	row := m.DB.QueryRowContext(ctx, query, user.FirstName, user.LastName, user.Email, hash)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (m *DBModel) UpdateUser(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		UPDATE 
			users 
		SET
		    first_name=$1,
			last_name=$2,
			email=$3
		WHERE id=$4
	`

	_, err := m.DB.ExecContext(ctx, query, user.FirstName, user.LastName, user.Email, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (m *DBModel) DeleteUser(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `DELETE FROM users WHERE id=$1`

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
