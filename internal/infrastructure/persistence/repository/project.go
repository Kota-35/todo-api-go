package repository

import (
	"context"
	"fmt"
	"todo-api-go/internal/domain/entity"
	"todo-api-go/internal/domain/repository"
	"todo-api-go/pkg/database"
	"todo-api-go/prisma/db"
)

type projectRepository struct {
	client *db.PrismaClient
}

func NewProjectRepository() repository.ProjectRepository {
	return &projectRepository{
		client: database.PrismaClient,
	}
}

func (p *projectRepository) Save(project *entity.Project) error {
	ctx := context.Background()

	if project.IsNew() {
		// 新規作成
		optionalData := []db.ProjectSetParam{}

		if description := project.Description(); description != nil {
			optionalData = append(
				optionalData,
				db.Project.Description.Set(*description),
			)
		}

		if color := project.Color(); color != nil {
			optionalData = append(optionalData, db.Project.Color.Set(*color))
		}

		createdProject, err := p.client.Project.CreateOne(
			db.Project.Name.Set(project.Name()),
			db.Project.Owner.Link(db.User.ID.Equals(project.OwnerID())),
			db.Project.Team.Link(db.Team.ID.Equals(project.TeamId())),
			optionalData...,
		).Exec(ctx)

		if err != nil {
			return fmt.Errorf("[projectRepository]プロジェクトの作成に失敗しました: %w", err)
		}

		return project.SetID(createdProject.ID)
	} else {
		// 更新
		updateData := []db.ProjectSetParam{
			db.Project.Name.Set(project.Name()),
		}

		if description := project.Description(); description != nil {
			updateData = append(updateData, db.Project.Description.Set(*description))
		}

		if color := project.Color(); color != nil {
			updateData = append(updateData, db.Project.Color.Set(*color))
		}

		_, err := p.client.Project.FindUnique(
			db.Project.ID.Equals(project.ID()),
		).Update(updateData...).Exec(ctx)

		if err != nil {
			return fmt.Errorf("[projectRepository]プロジェクトの更新に失敗しました: %w", err)
		}
	}

	return nil
}

func (p *projectRepository) FindByID(id string) (*entity.Project, error) {
	ctx := context.Background()

	project, err := p.client.Project.FindUnique(
		db.Project.ID.Equals(id),
	).Exec(ctx)

	if err != nil {
		return nil, fmt.Errorf("[projectRepository]プロジェクトの取得に失敗しました: %w", err)
	}

	reconstructProject := entity.ReconstructProject(
		project,
	)

	return &reconstructProject, nil
}

func (p *projectRepository) FindProjectsByTeamID(
	teamId string,
) ([]entity.Project, error) {
	ctx := context.Background()

	projects, err := p.client.Project.FindMany(
		db.Project.TeamID.Equals(teamId),
	).Exec(ctx)

	if err != nil {
		return nil, fmt.Errorf("[projectRepository]プロジェクトの取得に失敗しました: %w", err)
	}

	reconstructedProjects := entity.ReconstructProjects(projects)

	return reconstructedProjects, nil
}
