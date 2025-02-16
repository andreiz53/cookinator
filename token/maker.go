package token

import "time"

// Maker is an interface which manages tokens
type Maker interface {
	// CreateToken creates a new token for a specific user with the provided duration
	CreateToken(email string, duration time.Duration) (string, error)

	// VerifyToken checks the token validity
	VerifyToken(token string) (*Payload, error)
}
