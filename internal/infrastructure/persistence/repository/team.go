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
			return fmt.Errorf("[teamRepository]チームの作成に失敗しました: %w", err)
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
			return fmt.Errorf("[teamRepository.Save]チームの更新に失敗しました: %w", err)
		}
	}

	return nil
}

func (r *teamRepository) FindByID(id string) (*entity.Team, error) {
	ctx := context.Background()

	team, err := r.client.Team.FindUnique(
		db.Team.ID.Equals(id),
	).With(
		db.Team.Projects.Fetch(),
	).Exec(ctx)

	if err != nil {
		return nil, fmt.Errorf(
			"[teamRepository.FindByID]チームの取得に失敗しました: %w",
			err,
		)
	}

	reconstructTeam := entity.ReconstructTeam(team)

	return &reconstructTeam, nil
}

// ユーザーIDからチームを全て取得する
func (r *teamRepository) FindTeamsByUserID(
	userID string,
) ([]*entity.Team, error) {
	ctx := context.Background()

	user, err := r.client.User.FindUnique(
		db.User.ID.Equals(userID),
	).With(
		db.User.TeamMemberships.Fetch().With(
			db.TeamMember.Team.Fetch().With(
				db.Team.Projects.Fetch(),
			),
		),
		db.User.OwnedTeams.Fetch().With(
			db.Team.Projects.Fetch(),
		), // オーナーのチームも取得
	).Exec(ctx)

	if err != nil {
		return nil, fmt.Errorf(
			"[teamRepository.FindByUserID]ユーザーの取得に失敗しました: %w",
			err,
		)
	}

	teams := make([]*entity.Team, 0)
	// メンバーとして参加しているチーム
	for _, membership := range user.TeamMemberships() {
		team := membership.Team()
		reconstructTeam := entity.ReconstructTeam(team)
		teams = append(teams, &reconstructTeam)
	}

	for _, ownedTeam := range user.OwnedTeams() {
		reconstructTeam := entity.ReconstructTeam(&ownedTeam)
		teams = append(teams, &reconstructTeam)
	}

	return teams, nil
}
