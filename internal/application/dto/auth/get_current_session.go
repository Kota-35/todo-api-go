package auth

import "time"

type GetCurrentSessionInput struct {
	JWTToken string
}

type GetCurrentSessionOutput struct {
	User    GetCurrentSessionOutputUser  `json:"user"`
	Token   GetCurrentSessionOutputToken `json:"session"`
	Message string                       `json:"message"`
}

type GetCurrentSessionOutputUser struct {
	Id        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAT"`
}

type GetCurrentSessionOutputToken struct {
	ExpiresAt time.Time `json:"expiresAt"`
	IssuedAt  time.Time `json:"issuedAt"`
}
