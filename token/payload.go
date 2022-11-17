package token

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// NewPayload creates a new payload with the given username and duration
func NewPayload(username string, duration time.Duration) (*Payload, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	now := time.Now()

	return &Payload{
		ID:        id,
		Username:  username,
		IssuedAt:  now,
		ExpiresAt: now.Add(duration),
	}, nil
}

func (p *Payload) Valid() error {
	if p.ExpiresAt.Before(time.Now()) {
		return ErrExpiredToken
	}

	return nil
}
