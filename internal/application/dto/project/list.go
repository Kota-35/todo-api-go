package project

type GetProjectsOutput struct {
	Projects []GetProjectsOutputProject
}

type GetProjectsOutputProject struct {
	ProjectId   string  `json:"projectId"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Color       string  `json:"color"`
}
