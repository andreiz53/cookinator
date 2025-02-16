package util

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashed a given password using bcrypt with the default cost
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("failed to encrypt new password")
	}
	return string(hashedPassword), nil
}

// CheckPassword checks if the provided password matches with the hashed one
func CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
