package session_token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

func New() (plainToken string, tokenHash string, err error) {
	raw := make([]byte, 32)

	if _, err := rand.Read(raw); err != nil {
		return "", "", fmt.Errorf("generate session token: %w", err)
	}

	plainToken = base64.RawURLEncoding.EncodeToString(raw)
	tokenHash = Hash(plainToken)

	return plainToken, tokenHash, nil
}

func Hash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}
