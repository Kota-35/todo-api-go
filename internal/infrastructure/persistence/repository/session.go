package repository

import (
	"context"
	"fmt"
	"todo-api-go/internal/domain/entity"
	"todo-api-go/internal/domain/repository"
	valueobject "todo-api-go/internal/domain/valueobject"
	"todo-api-go/pkg/database"
	"todo-api-go/prisma/db"
)

type sessionRepository struct {
	client *db.PrismaClient
}

func NewSessionRepository() repository.SessionRepository {
	return &sessionRepository{
		client: database.PrismaClient,
	}
}

func (s *sessionRepository) Save(session *entity.Session) error {
	ctx := context.Background()

	if session.IsNew() {
		// 新規作成
		createdSession, err := s.client.Session.CreateOne(
			db.Session.TokenHash.Set(session.TokenHash()),
			db.Session.ExpiresAt.Set(session.ExpiresAt()),
			db.Session.User.Link(
				db.User.ID.Equals(session.UserId()),
			),
		).Exec(ctx)

		if err != nil {
			return fmt.Errorf("[sessionRepository]セッションの作成に失敗しました: %w", err)
		}

		return session.SetID(createdSession.ID)
	} else {
		return fmt.Errorf("[sessionRepository]セッションの更新はできません")
	}
}

func (s *sessionRepository) FindByToken(token *valueobject.RefreshToken) (*entity.Session, error) {
	ctx := context.Background()

	session, err := s.client.Session.FindUnique(
		db.Session.TokenHash.Equals(token.Hash()),
	).Exec(ctx)

	if err != nil {
		return nil, fmt.Errorf("[sessionRepository]セッションの検索に失敗しました: %w", err)
	}

	// domain entityに変換
	restoredSession := entity.ReconstructSession(
		session.ID,
		session.UserID,
		session.TokenHash,
		session.ExpiresAt,
		session.CreatedAt,
		session.IsRevoked,
	)

	return &restoredSession, nil
}

func (s *sessionRepository) FindByID(id string) (*entity.Session, error) {
	ctx := context.Background()

	session, err := s.client.Session.FindUnique(
		db.Session.ID.Equals(id),
	).Exec(ctx)

	if err != nil {
		return nil, fmt.Errorf("[sessionRepository]セッションの検索に失敗しました: %w", err)
	}

	// domain entityに変換
	restoredSession := entity.ReconstructSession(
		session.ID,
		session.UserID,
		session.TokenHash,
		session.ExpiresAt,
		session.CreatedAt,
		session.IsRevoked,
	)

	return &restoredSession, nil
}
