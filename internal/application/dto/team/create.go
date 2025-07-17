package team

import "time"

type CreateTeamInput struct {
	Name        string  `json:"Name"`
	Description *string `json:"description,omitempty"`
}

type CreateTeamOutput struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	OwnerId     string    `json:"ownerId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAT"`
}
