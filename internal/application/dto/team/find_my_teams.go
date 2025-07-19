package team

import "time"

type FindMyTeamsOutput struct {
	Teams []FindMyTeamsOutputTeam `json:"teams"`
}

type FindMyTeamsOutputTeam struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	OwnerId     string    `json:"ownerId"`
	CreatedAt   time.Time `json:"createdAt"`
}
