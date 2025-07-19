package team

import (
	"fmt"
	"todo-api-go/internal/domain/repository"

	teamDTO "todo-api-go/internal/application/dto/team"
)

type FindMyTeamsUseCase struct {
	teamRepo repository.TeamRepository
}

func NewFindMyTeamsUseCase(
	teamRepo repository.TeamRepository,
) *FindMyTeamsUseCase {
	return &FindMyTeamsUseCase{
		teamRepo: teamRepo,
	}
}

func (uc *FindMyTeamsUseCase) Exec(
	userId string,
) (*teamDTO.FindMyTeamsOutput, error) {
	// ユーザーIDから自分か所属するチームをすべて取得する
	teams, err := uc.teamRepo.FindTeamsByUserID(userId)
	if err != nil {
		return nil, fmt.Errorf("[FindMyTeamUseCase]複数チームの取得に失敗しました: %w", err)
	}

	outputTeams := make([]teamDTO.FindMyTeamsOutputTeam, len(teams))

	for i, team := range teams {
		outputTeams[i] = teamDTO.FindMyTeamsOutputTeam{
			Id:          team.ID(),
			Name:        team.Name(),
			Description: team.Description(),
			OwnerId:     team.OwnerID(),
			CreatedAt:   team.CreatedAt(),
		}
	}

	return &teamDTO.FindMyTeamsOutput{
		Teams: outputTeams,
	}, nil

}
