package auth

import (
	"time"
	"todo-api-go/internal/domain/valueobject"
)

type RefreshSessionInput struct {
	RefreshTokenVO valueobject.RefreshToken
}

type RefreshSessionOutput struct {
	AccessToken          string
	AccessTokenExpiresAt time.Time
}
