package security

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTGenerator interface {
	GenerateAccessToken(uid string, ssID string, expiresAt time.Time) (string, error)
}

type jwtGenerator struct {
	jwtSecret string
}

func NewJWTGenerator(jwtSecret string) JWTGenerator {
	return &jwtGenerator{
		jwtSecret: jwtSecret,
	}
}

func (g *jwtGenerator) GenerateAccessToken(uid string, rtID string, expiresAt time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":          uid,
		"refresh_token_id": rtID,
		"expires_at":       expiresAt.Unix(),
	})

	tokenString, err := token.SignedString([]byte(g.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("[jwtGenerator]JWTの生成に失敗しました: %w", err)
	}

	return tokenString, nil
}
