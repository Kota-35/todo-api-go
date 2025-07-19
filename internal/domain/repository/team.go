package repository

import "todo-api-go/internal/domain/entity"

type TeamRepository interface {
	Save(team *entity.Team) error
	FindByID(id string) (*entity.Team, error)
	FindTeamsByUserID(userID string) ([]*entity.Team, error)
}
