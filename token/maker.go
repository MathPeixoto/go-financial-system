package token

import "time"

// Maker is an interface for managing tokens.
type Maker interface {
	// CreateToken creates a token for a specific username and duration
	CreateToken(username string, duration time.Duration) (string, *Payload, error)

	// VerifyToken verifies a token and returns the payload associated with it
	VerifyToken(token string) (*Payload, error)
}
