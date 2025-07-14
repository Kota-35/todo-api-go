package auth

import "time"

type AuthenticateUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthenticateUserOutput struct {
	AccessToken           string
	AccessTokenExpiresAt  time.Time
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
}
