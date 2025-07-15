package auth

import (
	"fmt"
	authDTO "todo-api-go/internal/application/dto/auth"
	"todo-api-go/internal/domain/repository"
	"todo-api-go/internal/domain/valueobject"
)

type LogoutUseCase struct {
	pepper      valueobject.Pepper
	sessionRepo repository.SessionRepository
}

func NewLogoutUserCase(
	pepper valueobject.Pepper,
	sessionRepo repository.SessionRepository,
) *LogoutUseCase {
	return &LogoutUseCase{
		pepper:      pepper,
		sessionRepo: sessionRepo,
	}
}

func (uc *LogoutUseCase) Execute(input authDTO.LogoutInput) error {

	// sessionの取得
	session, err := uc.sessionRepo.FindByID(input.RefreshTokenId)
	if err != nil {
		return fmt.Errorf("[LogoutUseCase]セッションの取得に失敗しました: %w", err)
	}

	session.Revoke()

	// fmt.Sprint("Session: %w", session)

	if err := uc.sessionRepo.Save(session); err != nil {
		return fmt.Errorf("[LogoutUseCase]セッションを更新に失敗しました: %w", err)
	}

	return nil
}
