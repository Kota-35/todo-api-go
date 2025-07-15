package auth

import (
	"fmt"
	"todo-api-go/internal/domain/repository"
	"todo-api-go/internal/domain/security"

	authDTO "todo-api-go/internal/application/dto/auth"
)

type GetCurrentSessionUseCase struct {
	jwtGenerator security.JWTGenerator
	userRepo     repository.UserRepository
	sessionRepo  repository.SessionRepository
}

func NewGetCurrentSessionUseCase(
	jwtGenerator security.JWTGenerator,
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
) *GetCurrentSessionUseCase {
	return &GetCurrentSessionUseCase{
		jwtGenerator: jwtGenerator,
		userRepo:     userRepo,
		sessionRepo:  sessionRepo,
	}
}

func (uc *GetCurrentSessionUseCase) Execute(
	input authDTO.GetCurrentSessionInput,
) (*authDTO.GetCurrentSessionOutput, error) {
	// アクセストークンの検証
	claims, err := uc.jwtGenerator.VerifyAccessToken(input.JWTToken)

	if err != nil {
		return nil, fmt.Errorf(
			"[GetCurrentSessionUseCase]トークンの検証に失敗しました: %w",
			err,
		)
	}

	session, err := uc.sessionRepo.FindByID(claims.RefreshTokenID)
	if err != nil {
		return nil, fmt.Errorf(
			"[GetCurrentSessionUseCase]セッションの取得に失敗しました:\n %w", err,
		)
	}

	if session.IsRevoked() {
		return nil, fmt.Errorf("[GetCurrentSessionUseCase]無効なセッションです")
	}

	user, err := uc.userRepo.FindByID(claims.UserID)

	if err != nil {
		return nil, fmt.Errorf(
			"[GetCurrentSessionUseCase]ユーザーの取得に失敗しました: %w",
			err,
		)
	}

	return &authDTO.GetCurrentSessionOutput{
		User: authDTO.GetCurrentSessionOutputUser{
			Id:        user.ID(),
			Username:  user.Username(),
			Email:     user.Email().String(),
			IsActive:  user.IsActive(),
			CreatedAt: user.CreatedAt(),
		},
		Token: authDTO.GetCurrentSessionOutputToken{
			ExpiresAt: claims.ExpiresAt,
			IssuedAt:  claims.Iat,
		},
	}, nil

}
