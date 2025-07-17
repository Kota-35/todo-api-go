package entity

import (
	"errors"
	"time"
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
	id string,
	name string,
	description *string,
	ownerId string,
	createdAt time.Time,
	updatedAt time.Time,
	projects []Project,
) Team {
	return Team{
		id:          id,
		name:        name,
		description: description,
		ownerId:     ownerId,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
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
