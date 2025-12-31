package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/fx"
	"soft.structx.io/dino/setup"
)

// Claims
type Claims struct {
	jwt.RegisteredClaims
}

// Authenticator
type Authenticator interface {
	GenerateJWT(string, string, string) (string, error)
}

// Params
type Params struct {
	fx.In

	Cfg *setup.Authenticator
}

// Result
type Result struct {
	fx.Out

	Auth Authenticator
}

type simpleAuth struct {
	issuer      string
	audience    []string
	jwtDuration time.Duration
}

// interface compliance
var _ Authenticator = (*simpleAuth)(nil)

// Module
var Module = fx.Module("authenticator", fx.Provide(newModule))

func newModule(p Params) Result {
	return Result{
		Auth: &simpleAuth{
			issuer:      p.Cfg.JWT.Issuer,
			audience:    p.Cfg.JWT.Audience,
			jwtDuration: time.Duration(p.Cfg.JWT.Duration),
		},
	}
}

// GenerateJWT implements Authenticator.
func (s *simpleAuth) GenerateJWT(subject, id, signingKey string) (string, error) {
	c := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        id,
			Issuer:    s.issuer,
			Subject:   subject,
			Audience:  []string{"dino"},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			// ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.jwtDuration)),
		},
	}

	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	ss, err := tk.SignedString([]byte(signingKey))
	if err != nil {
		return "", fmt.Errorf("token.SignedString: %w", err)
	}

	return ss, nil
}
