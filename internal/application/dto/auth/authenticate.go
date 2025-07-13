package auth

import "time"

type AuthenticateUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthenticateUserOutput struct {
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	ExpiresAt    time.Time `json:"expiresAt"`
}
