package util

import (
	"crypto/rand"
	"fmt"
)

func GenerateShortID(length int, exists func(string) bool) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const maxAttempts = 100

	for i := 0; i < maxAttempts; i++ {
		b := make([]byte, length)
		if _, err := rand.Read(b); err != nil {
			return "", err
		}

		for j := range b {
			b[j] = charset[b[j]%byte(len(charset))]
		}

		id := string(b)
		if !exists(id) {
			return id, nil
		}
	}

	return "", fmt.Errorf("failed to generate unique ID after %d attempts", maxAttempts)
}
