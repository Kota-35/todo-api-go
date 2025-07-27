package project

import (
	"fmt"
	projectDTO "todo-api-go/internal/application/dto/project"
	"todo-api-go/internal/domain/entity"
	"todo-api-go/internal/domain/repository"
)

type CreateProjectUseCase struct {
	projectRepository repository.ProjectRepository
}

func NewCreateProjectUseCase(
	projectRepository repository.ProjectRepository,
) *CreateProjectUseCase {
	return &CreateProjectUseCase{
		projectRepository: projectRepository,
	}
}

func (uc *CreateProjectUseCase) Execute(
	input projectDTO.CreateProjectInput,
	userId *string,
	teamId string,
) (*projectDTO.CreateProjectOutput, error) {

	// エンティティの作成
	project := entity.NewProject(
		input.Name,
		input.Description,
		input.Color,
		*userId,
		teamId,
	)

	// 永続化
	if err := uc.projectRepository.Save(&project); err != nil {
		return nil, fmt.Errorf(
			"[CreateProjectUseCase]プロジェクトの作成に失敗しました: %w",
			err,
		)
	}

	return &projectDTO.CreateProjectOutput{
		ProjectId:   project.ID(),
		Name:        project.Name(),
		Description: project.Description(),
		Color:       project.Color(),
	}, nil
}
