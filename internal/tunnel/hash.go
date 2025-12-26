package tunnel

import (
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"
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

func verifyToken(token string, combinedHash string) (bool, error) {
	parts := strings.Split(combinedHash, ":")
	if len(parts) != 2 {
		return false, fmt.Errorf("stored hash format is invalid expected (hash:salt)")
	}

	encodedHash, encodedSalt := parts[0], parts[1]

	storedHash, err := base64.StdEncoding.DecodeString(encodedHash)
	if err != nil {
		return false, fmt.Errorf("base64 decode hash: %w", err)
	}

	storedSalt, err := base64.StdEncoding.DecodeString(encodedSalt)
	if err != nil {
		return false, fmt.Errorf("base64 decode salt: %w", err)
	}

	derivedKey, err := pbkdf2.Key(
		sha512.New,
		token,
		storedSalt,
		iterations,
		keyLength,
	)
	if err != nil {
		return false, fmt.Errorf("pbkdf2.Key: %w", err)
	}

	return subtle.ConstantTimeCompare(derivedKey, storedHash) == 1, nil
}
