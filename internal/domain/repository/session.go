package repository

import (
	"todo-api-go/internal/domain/entity"
	"todo-api-go/internal/domain/valueobject"
)

type SessionRepository interface {
	Save(session *entity.Session) error
	FindByToken(h *valueobject.RefreshToken) (*entity.Session, error)
	FindByID(id string) (*entity.Session, error)
}
