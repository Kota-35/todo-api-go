package auth

import (
	"fmt"
	"time"
	"todo-api-go/internal/domain/entity"
	"todo-api-go/internal/domain/repository"
	"todo-api-go/internal/domain/security"
	"todo-api-go/internal/domain/valueobject"

	authDTO "todo-api-go/internal/application/dto/auth"
	domainError "todo-api-go/internal/domain/error"
)

type AuthenticateUserUseCase struct {
	userRepo        repository.UserRepository
	sessionRepo     repository.SessionRepository
	jwtGenerator    security.JWTGenerator
	refreshTokenGen security.RefreshTokenGenerator
	pepper          valueobject.Pepper
	sessionTTL      time.Duration
}

func NewAuthenticateUserUseCase(
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
	jwtGenerator security.JWTGenerator,
	refreshTokenGen security.RefreshTokenGenerator,
	pepper valueobject.Pepper,
	sessionTTL time.Duration,
) *AuthenticateUserUseCase {
	return &AuthenticateUserUseCase{
		userRepo:        userRepo,
		sessionRepo:     sessionRepo,
		jwtGenerator:    jwtGenerator,
		refreshTokenGen: refreshTokenGen,
		pepper:          pepper,
		sessionTTL:      sessionTTL,
	}
}

func (uc *AuthenticateUserUseCase) Execute(input authDTO.AuthenticateUserInput) (*authDTO.AuthenticateUserOutput, error) {
	// ユーザーの取得
	emailVO, err := valueobject.NewEmail(input.Email)
	if err != nil {
		return nil, domainError.NewInvalidUserDataError(
			fmt.Sprintf("メールアドレスに問題があります: %s", err.Error()),
		)
	}

	user, err := uc.userRepo.FindByEmail(emailVO)
	if err != nil {
		return nil, domainError.NewInvalidUserDataError(
			fmt.Sprintf("メールアドレスまたはパスワードが正しくありません: %s", err.Error()),
		)
	}

	if !user.VerifyPassword(input.Password) {
		return nil, domainError.NewAuthenticationError(
			"パスワードが正しくありません",
		)
	}

	// リフレッシュトークンの生成（ドメインサービスを使用）
	refreshToken, err := uc.refreshTokenGen.Generate()
	if err != nil {
		return nil, fmt.Errorf("リフレッシュトークンの生成に失敗しました: %w", err)
	}

	// セッションの作成
	session, err := entity.NewSession(
		user.ID(),
		refreshToken,
		uc.pepper,
		uc.sessionTTL,
	)
	if err != nil {
		return nil, fmt.Errorf("セッションの作成に失敗しました: %w", err)
	}

	// セッションの永続化
	if err := uc.sessionRepo.Save(&session); err != nil {
		return nil, fmt.Errorf("セッションの永続化に失敗しました: %w", err)
	}

	// アクセストークンの生成
	accessToken, err := uc.jwtGenerator.GenerateAccessToken(
		user.ID(),
		session.ID(),
		session.ExpiresAt(),
	)
	if err != nil {
		return nil, fmt.Errorf("アクセストークンの生成に失敗しました: %w", err)
	}

	// レスポンスの作成
	return &authDTO.AuthenticateUserOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken, // 平文のトークンを返す
		ExpiresAt:    session.ExpiresAt(),
	}, nil
}
