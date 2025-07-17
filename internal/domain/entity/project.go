package entity

import (
	"errors"
	"time"
)

type Project struct {
	id          string
	name        string
	description *string
	color       *string
	ownerId     string
	teamId      *string
	createdAt   time.Time
	updatedAt   time.Time
}

func NewProject(
	id string,
	name string,
	description *string,
	color *string,
	ownedId string,
	teamId *string,
	createdAt time.Time,
	updatedAt time.Time,
) Project {
	return Project{
		id:          "",
		name:        name,
		description: description,
		color:       color,
		ownerId:     ownedId,
		teamId:      teamId,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

func ReconstructProject(
	id string,
	name string,
	description *string,
	color *string,
	ownedId string,
	teamId *string,
	createdAt time.Time,
	updatedAt time.Time,
) Project {
	return Project{
		id:          id,
		name:        name,
		description: description,
		color:       color,
		ownerId:     ownedId,
		teamId:      teamId,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

func (p *Project) IsNew() bool {
	return p.id == ""
}

func (p *Project) Name() string {
	return p.name
}

func (p *Project) Description() *string {
	return p.description
}

func (p *Project) Color() *string {
	return p.color
}

func (p *Project) OwnerID() string {
	return p.ownerId
}

func (p *Project) TeamId() *string {
	return p.teamId
}

func (p *Project) CreatedAt() time.Time {
	return p.createdAt
}

func (p *Project) SetID(id string) error {
	if p.id != "" {
		return errors.New("IDはすでに設定されています")
	}
	if id == "" {
		return errors.New("IDは空にできません")
	}

	p.id = id
	return nil
}
