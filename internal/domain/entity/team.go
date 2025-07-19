package entity

import (
	"errors"
	"time"

	"todo-api-go/prisma/db"
)

type Team struct {
	id          string
	name        string
	description *string
	ownerId     string
	createdAt   time.Time
	updatedAt   time.Time

	projects []Project
}

func NewTeam(
	name string,
	description *string,
	ownerId string,
) Team {
	return Team{
		id:          "",
		name:        name,
		description: description,
		ownerId:     ownerId,
		createdAt:   time.Now(),
		updatedAt:   time.Now(),
	}
}

func ReconstructTeam(
	teamModel *db.TeamModel,
) Team {

	var descPtr *string
	description, ok := teamModel.Description()
	if ok {
		descPtr = &description
	}

	projects := ReconstructProjects(teamModel.Projects())

	return Team{
		id:          teamModel.ID,
		name:        teamModel.Name,
		description: descPtr,
		ownerId:     teamModel.OwnerID,
		createdAt:   teamModel.CreatedAt,
		updatedAt:   teamModel.UpdatedAt,
		projects:    projects,
	}
}

func (t *Team) IsNew() bool {
	return t.id == ""
}

func (t *Team) ID() string {
	return t.id
}

func (t *Team) Name() string {
	return t.name
}

func (t *Team) Description() *string {
	return t.description
}

func (t *Team) OwnerID() string {
	return t.ownerId
}

func (t *Team) Projects() []Project {
	return t.projects
}

func (t *Team) CreatedAt() time.Time {
	return t.createdAt
}

func (t *Team) UpdatedAt() time.Time {
	return t.updatedAt
}

func (t *Team) SetID(id string) error {
	if t.id != "" {
		return errors.New("IDはすでに設定されています")
	}
	if id == "" {
		return errors.New("IDは空にできません")
	}

	t.id = id
	return nil
}
