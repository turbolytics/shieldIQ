package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateSecret() (string, error) {
	// Generate a random 32-byte secret
	secretBytes := make([]byte, 32)
	_, err := rand.Read(secretBytes)
	if err != nil {
		return "", err
	}

	// Encode the bytes to a hexadecimal string
	return hex.EncodeToString(secretBytes), nil
}
