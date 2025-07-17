package team

import (
	"fmt"
	"todo-api-go/internal/domain/entity"
	"todo-api-go/internal/domain/repository"

	teamDTO "todo-api-go/internal/application/dto/team"
)

type CreateTeamUseCase struct {
	teamRepository repository.TeamRepository
}

func NewCreateTeamUseCase(
	teamRepository repository.TeamRepository,
) *CreateTeamUseCase {
	return &CreateTeamUseCase{
		teamRepository: teamRepository,
	}
}

func (uc *CreateTeamUseCase) Execute(
	input teamDTO.CreateTeamInput,
	userId string,
) (*teamDTO.CreateTeamOutput, error) {

	// エンティティの作成
	team := entity.NewTeam(
		input.Name,
		input.Description,
		userId,
	)

	// 永続化
	if err := uc.teamRepository.Save(&team); err != nil {
		return nil, fmt.Errorf("[CreateTeamUseCase]チームの永続化に失敗しました: %w", err)
	}

	// INFO: Descriptionはoptionalのため、更新でしかセットできない
	if input.Description != nil {
		if err := uc.teamRepository.Save(&team); err != nil {
			return nil, fmt.Errorf("[CreateTeamUseCase]チームの更新に失敗しました: %w", err)
		}
	}

	return &teamDTO.CreateTeamOutput{
		Id:          team.ID(),
		Name:        team.Name(),
		Description: team.Description(),
		OwnerId:     team.OwnerID(),
		CreatedAt:   team.CreatedAt(),
		UpdatedAt:   team.UpdatedAt(),
	}, nil
}
