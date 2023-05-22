package models

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
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
