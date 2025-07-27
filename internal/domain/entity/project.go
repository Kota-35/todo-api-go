package entity

import (
	"errors"
	"time"
	"todo-api-go/prisma/db"
)

type Project struct {
	id          string
	name        string
	description *string
	color       *string
	ownerId     string
	teamId      string
	createdAt   time.Time
	updatedAt   time.Time
}

func NewProject(
	name string,
	description *string,
	color *string,
	ownedId string,
	teamId string,
) Project {
	return Project{
		id:          "",
		name:        name,
		description: description,
		color:       color,
		ownerId:     ownedId,
		teamId:      teamId,
		createdAt:   time.Now(),
		updatedAt:   time.Now(),
	}
}

func ReconstructProject(
	projectModel *db.ProjectModel,
) Project {
	var descPtr *string
	description, ok := projectModel.Description()
	if ok {
		descPtr = &description
	}

	return Project{
		id:          projectModel.ID,
		name:        projectModel.Name,
		description: descPtr,
		color:       &projectModel.Color,
		ownerId:     projectModel.OwnerID,
		teamId:      projectModel.TeamID,
		createdAt:   projectModel.CreatedAt,
		updatedAt:   projectModel.UpdatedAt,
	}
}

// 　プロジェクトモデルの配列からプロジェクトの配列に変換
func ReconstructProjects(
	projectModels []db.ProjectModel,
) []Project {
	projects := make([]Project, len(projectModels))
	for i, v := range projectModels {
		projects[i] = ReconstructProject(&v)
	}

	return projects
}

func (p *Project) IsNew() bool {
	return p.id == ""
}

func (p *Project) Name() string {
	return p.name
}

func (p *Project) Description() *string {
	return p.description
}

func (p *Project) Color() *string {
	return p.color
}

func (p *Project) OwnerID() string {
	return p.ownerId
}

func (p *Project) TeamId() string {
	return p.teamId
}

func (p *Project) CreatedAt() time.Time {
	return p.createdAt
}

func (p *Project) SetID(id string) error {
	if p.id != "" {
		return errors.New("IDはすでに設定されています")
	}
	if id == "" {
		return errors.New("IDは空にできません")
	}

	p.id = id
	return nil
}

func (p *Project) ID() string {
	return p.id
}
