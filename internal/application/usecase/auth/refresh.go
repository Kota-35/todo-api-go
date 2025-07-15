package auth

import (
	"fmt"

	"time"
	authDTO "todo-api-go/internal/application/dto/auth"
	"todo-api-go/internal/domain/repository"
	"todo-api-go/internal/domain/security"
)

type RefreshSessionUseCase struct {
	sessionRepo repository.SessionRepository
	userRepo    repository.UserRepository

	jwtGenerator security.JWTGenerator
}

func NewRefreshSessionUseCase(
	sessionRepo repository.SessionRepository,
	userRepo repository.UserRepository,
	jwtGenerator security.JWTGenerator,
) *RefreshSessionUseCase {
	return &RefreshSessionUseCase{
		sessionRepo: sessionRepo,
		userRepo:    userRepo,

		jwtGenerator: jwtGenerator,
	}
}

func (uc *RefreshSessionUseCase) Execute(
	input *authDTO.RefreshSessionInput,
) (*authDTO.RefreshSessionOutput, error) {
	// セッションの検索と検証
	session, err := uc.sessionRepo.FindByToken(&input.RefreshTokenVO)

	if err != nil {
		return nil, fmt.Errorf(
			"[RefreshSessionUseCase]セッションの取得に失敗しました: %w",
			err,
		)
	}

	// セッションの有効性検証
	// ユーザーがログアウトしたり、セキュリティの理由で明示的に無効化されたセッション
	if session.IsRevoked() {
		return nil, fmt.Errorf("[RefreshSessionUseCase]セッションが無効です")
	}

	// セッションの有効期限検証
	// 時間経過により自動的に無効になったセッション
	if !session.IsActive(time.Now()) {
		return nil, fmt.Errorf("[RefreshSessionUseCase]セッションの有効期限が切れています")
	}

	// ユーザーの検証
	user, err := uc.userRepo.FindByID(session.UserId())
	if err != nil {
		return nil, fmt.Errorf(
			"[RefreshSessionUseCase]ユーザーの取得に失敗しました: %w",
			err,
		)
	}

	if !user.IsActive() {
		return nil, fmt.Errorf("[RefreshSessionUseCase]非アクティブなユーザーです: %w", err)
	}

	accessTokenExpiresAt := time.Now().Add(AccessTokenTTL)

	// アクセストークンの生成
	accessToken, err := uc.jwtGenerator.GenerateAccessToken(
		user.ID(),
		session.ID(),
		accessTokenExpiresAt,
	)
	if err != nil {
		return nil, fmt.Errorf("アクセストークンの生成に失敗しました: %w", err)
	}

	return &authDTO.RefreshSessionOutput{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessTokenExpiresAt,
	}, nil

}
