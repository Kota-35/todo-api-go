package repository

import (
	"context"
	"fmt"
	"todo-api-go/internal/domain/entity"
	"todo-api-go/internal/domain/repository"
	"todo-api-go/pkg/database"
	"todo-api-go/prisma/db"
)

type teamRepository struct {
	client *db.PrismaClient
}

func NewTeamRepository() repository.TeamRepository {
	return &teamRepository{
		client: database.PrismaClient,
	}
}

func (r *teamRepository) Save(team *entity.Team) error {
	ctx := context.Background()

	if team.IsNew() {
		// 新規作成
		createdTeam, err := r.client.Team.CreateOne(
			db.Team.Name.Set(team.Name()),
			db.Team.Owner.Link(db.User.ID.Equals(team.OwnerID())),
		).Exec(ctx)

		if err != nil {
			return fmt.Errorf("[userRepository]ユーザーの作成に失敗しました: %w", err)
		}

		return team.SetID(createdTeam.ID)
	} else {
		// 更新
		_, err := r.client.Team.FindUnique(
			db.Team.ID.Equals(team.ID()),
		).Update(
			db.Team.Name.Set(team.Name()),
			db.Team.Description.Set(*team.Description()),
			db.Team.Owner.Link(db.User.ID.Equals(team.OwnerID())),
		).Exec(ctx)

		if err != nil {
			return fmt.Errorf("[teamRepository.Save]ユーザーの更新に失敗しました: %w", err)
		}
	}

	return nil
}

func (r *teamRepository) FindByID(id string) (*entity.Team, error) {
	ctx := context.Background()

	team, err := r.client.Team.FindUnique(
		db.Team.ID.Equals(id),
	).Exec(ctx)

	if err != nil {
		return nil, fmt.Errorf(
			"[teamRepository.FindByID]チームの取得に失敗しました: %w",
			err,
		)
	}

	description, ok := team.Description()
	var descPtr *string
	if ok {
		descPtr = &description
	}

	tempProjects := make([]entity.Project, len(team.Projects()))
	for i, v := range team.Projects() {
		description, ok := team.Description()
		var tempDescPtr *string
		if ok {
			tempDescPtr = &description
		}
		tempProjects[i] = entity.ReconstructProject(
			v.ID,
			v.Name,
			tempDescPtr,
			&v.Color,
			v.OwnerID,
			&v.TeamID,
			v.CreatedAt,
			v.UpdatedAt,
		)
	}

	reconstructTeam := entity.ReconstructTeam(
		team.ID,
		team.Name,
		descPtr,
		team.OwnerID,
		team.CreatedAt,
		team.UpdatedAt,
		tempProjects,
	)

	return &reconstructTeam, nil
}
