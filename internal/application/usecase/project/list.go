package project

import (
	"fmt"
	"todo-api-go/internal/domain/repository"

	projectDTO "todo-api-go/internal/application/dto/project"
)

type GetProjectsUseCase struct {
	projectRepository repository.ProjectRepository
}

func NewGetProjectsUseCase(
	projectRepository repository.ProjectRepository,
) *GetProjectsUseCase {
	return &GetProjectsUseCase{
		projectRepository: projectRepository,
	}
}

func (uc *GetProjectsUseCase) Execute(
	teamId string,
) (*projectDTO.GetProjectsOutput, error) {

	// プロジェクトの一覧を取得
	projects, err := uc.projectRepository.FindProjectsByTeamID(teamId)

	if err != nil {
		return nil, fmt.Errorf("[GetProjectsUseCase]プロジェクトの取得に失敗しました: %w", err)
	}

	outputProjects := make(
		[]projectDTO.GetProjectsOutputProject,
		len(projects),
	)

	for i, project := range projects {
		outputProjects[i] = projectDTO.GetProjectsOutputProject{
			ProjectId:   project.ID(),
			Name:        project.Name(),
			Description: project.Description(),
			Color:       *project.Color(),
		}
	}

	return &projectDTO.GetProjectsOutput{
		Projects: outputProjects,
	}, nil
}
