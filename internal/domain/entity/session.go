package entity

import (
	"errors"
	"time"
	"todo-api-go/internal/domain/valueobject"
)

type Session struct {
	id        string
	userId    string
	tokenHash string
	expiresAt time.Time
	createdAt time.Time
	isRevoked bool
}

func NewSession(
	userId string,
	token string,
	pepper valueobject.Pepper,
	ttl time.Duration,
) (Session, error) {
	refreshTokenVO, err := valueobject.NewRefreshToken(token, pepper)
	if err != nil {
		return Session{}, err
	}

	now := time.Now()

	return Session{
		id:        "",
		userId:    userId,
		tokenHash: refreshTokenVO.Hash(),
		expiresAt: now.Add(ttl),
		createdAt: time.Now(),
		isRevoked: false,
	}, nil
}

// データベースから復元されたセッションを作成
func ReconstructSession(
	id, userId, tokenHash string,
	expiresAt, createdAt time.Time,
	isRevoked bool,
) Session {
	return Session{
		id:        id,
		userId:    userId,
		tokenHash: tokenHash,
		expiresAt: expiresAt,
		createdAt: createdAt,
		isRevoked: isRevoked,
	}
}

// 状態遷移
func (s *Session) Revoke() { s.isRevoked = true }
func (s Session) IsActive(now time.Time) bool {
	return now.Before(s.expiresAt)
}

func (s *Session) IsNew() bool {
	return s.id == ""
}

func (s *Session) TokenHash() string {
	return s.tokenHash
}

func (s Session) UserId() string {
	return s.userId
}

func (s *Session) ExpiresAt() time.Time {
	return s.expiresAt
}

func (s *Session) SetID(id string) error {
	if s.id != "" {
		return errors.New("IDはすでに設定されています")
	}

	if id == "" {
		return errors.New("IDは空にできません")
	}

	s.id = id
	return nil
}

func (s *Session) ID() string {
	return s.id
}
