package token

import (
	"fmt"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20"
	"golang.org/x/crypto/chacha20poly1305"
	"time"
)

// PasetoMaker is a Maker implementation that uses Paseto for token creation and verification.
type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

// NewPasetoMaker creates a new PasetoMaker.
func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("symmetric key must be %d bytes", chacha20.KeySize)
	}

	return &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}, nil
}

// CreateToken creates a token for a specific username and duration.
func (m *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	return m.paseto.Encrypt(m.symmetricKey, payload, nil)
}

// VerifyToken verifies a token and returns the payload associated with it.
func (m *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	if err := m.paseto.Decrypt(token, m.symmetricKey, payload, nil); err != nil {
		return nil, ErrInvalidToken
	}

	err := payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
