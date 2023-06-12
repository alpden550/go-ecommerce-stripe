package models

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	ScopeAuthentication = "authentication"
)

type Token struct {
	PlainText   string    `json:"plain_text"`
	UserID      int64     `json:"-"`
	Hash        []byte    `json:"-"`
	ExpiredDate time.Time `json:"expired_date"`
	Scope       string    `json:"-"`
}

func GenerateNewToken(userID int, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID:      int64(userID),
		ExpiredDate: time.Now().Add(ttl),
		Scope:       scope,
	}

	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.PlainText))
	token.Hash = hash[:]

	return token, nil
}

func (m *DBModel) InsertToken(token *Token, user *User) error {
	deleteQuery := `DELETE FROM tokens WHERE user_id=$1`
	_, err := m.DB.Exec(deleteQuery, user.ID)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO tokens (user_id, name, email, hash, expire_date)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err = m.DB.Exec(query, user.ID, user.LastName, user.Email, token.Hash, token.ExpiredDate)
	if err != nil {
		return err
	}

	return nil
}

func (m *DBModel) GetUserForToken(token string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tokenHash := sha256.Sum256([]byte(token))
	var user User
	query := `SELECT 
    			u.id, u.first_name, u.last_name, u.email 
			  FROM users u
			  INNER JOIN tokens t on u.id = t.user_id
			  WHERE t.hash=$1 AND t.expire_date > $2
	`

	if err := m.DB.QueryRowContext(ctx, query, tokenHash[:], time.Now()).
		Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
		); err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *DBModel) Authenticate(email, password string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string
	query := `SELECT id, password from users WHERE email=$1`

	if err := m.DB.QueryRowContext(ctx, query, email).Scan(&id, &hashedPassword); err != nil {
		return id, err
	}

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return 0, errors.New("incorrect password")
	} else if err != nil {
		return 0, fmt.Errorf("%w", err)
	}

	return id, nil
}

func (m *DBModel) DeleteToken(userID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `DELETE FROM tokens WHERE user_id=$1`
	_, err := m.DB.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil
}
