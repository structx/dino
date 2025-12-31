package tunnel

import (
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
)

const (
	iterations = 210000
	saltLength = 16
	keyLength  = 64
)

func hashToken(token string) (string, error) {
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("rand.Read: %w", err)
	}

	derivedKey, err := pbkdf2.Key(
		sha512.New,
		token,
		salt,
		iterations,
		keyLength,
	)
	if err != nil {
		return "", fmt.Errorf("pbkdf2.New: %w", err)
	}

	encodedHash := base64.StdEncoding.EncodeToString(derivedKey)
	encodedSalt := base64.StdEncoding.EncodeToString(salt)

	return fmt.Sprintf("%s:%s", encodedHash, encodedSalt), nil
}
