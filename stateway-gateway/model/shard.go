package model

import (
	"time"

	"github.com/disgoorg/snowflake/v2"
	"gopkg.in/guregu/null.v4"
)

type ShardSession struct {
	ID            string
	AppID         snowflake.ID
	ShardID       int
	LastSequence  int
	ResumeURL     string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	InvalidatedAt null.Time
}
