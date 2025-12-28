package verifier

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"soft.structx.io/dino/auth"
	"soft.structx.io/dino/internal/tunnel"
)

// Verifier
type Verifier interface {
	VerifyToken(context.Context, string, string) (*auth.Claims, error)
}

type jwtVerifier struct {
	t tunnel.Service
}

func newVerifier(tunelSvc tunnel.Service) Verifier {
	return &jwtVerifier{
		t: tunelSvc,
	}
}

// VerifyToken implements Authenticator.
func (j *jwtVerifier) VerifyToken(ctx context.Context, tunnelID, token string) (*auth.Claims, error) {
	combinedHash, err := j.t.VerifyToken(ctx, tunnelID)
	if err != nil {
		return nil, fmt.Errorf("tunnelSvc.VerifyToken: %w", err)
	}

	if tk, err := jwt.ParseWithClaims(token, &auth.Claims{}, func(t *jwt.Token) (any, error) {
		return []byte(combinedHash), nil
	}); err != nil {
		return nil, fmt.Errorf("jwt.ParseWithClaims: %w", err)
	} else if claims, ok := tk.Claims.(*auth.Claims); ok {
		return claims, nil
	} else {
		return nil, errors.New("unknown claims")
	}
}
