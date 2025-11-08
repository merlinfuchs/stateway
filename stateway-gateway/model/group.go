package model

import "time"

type Group struct {
	ID                 string
	DisplayName        string
	DefaultConstraints AppConstraints
	DefaultConfig      AppConfig
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
