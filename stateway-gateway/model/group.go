package model

import "time"

type Group struct {
	ID          string
	DisplayName string
	Constraints AppConstraints
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
