package repository

import "todo-api-go/internal/domain/entity"

type ProjectRepository interface {
	Save(project *entity.Project) error
	FindByID(id string) (*entity.Project, error)
	// FindByOwnerID(ownerId string) (*entity.Project, error)
	FindProjectsByTeamID(teamId string) ([]entity.Project, error)
}
