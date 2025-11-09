package model

import (
	"time"

	"github.com/disgoorg/snowflake/v2"
	"gopkg.in/guregu/null.v4"
)

type AppDisabledCode string

const (
	AppDisabledCodeUnknown           AppDisabledCode = "unknown"
	AppDisabledCodeInvalidToken      AppDisabledCode = "invalid_token"
	AppDisabledCodeInvalidIntents    AppDisabledCode = "invalid_intents"
	AppDisabledCodeDisallowedIntents AppDisabledCode = "disallowed_intents"
)

type App struct {
	ID                  snowflake.ID    `json:"id"`
	GroupID             string          `json:"group_id"`
	DisplayName         string          `json:"display_name"`
	DiscordClientID     snowflake.ID    `json:"discord_client_id"`
	DiscordBotToken     string          `json:"discord_bot_token"`
	DiscordPublicKey    string          `json:"discord_public_key"`
	DiscordClientSecret null.String     `json:"discord_client_secret"`
	ShardCount          int             `json:"shard_count"`
	Disabled            bool            `json:"disabled"`
	DisabledCode        AppDisabledCode `json:"disabled_code"`
	DisabledMessage     null.String     `json:"disabled_message"`
	Constraints         AppConstraints  `json:"constraints"`
	Config              AppConfig       `json:"config"`
	CreatedAt           time.Time       `json:"created_at"`
	UpdatedAt           time.Time       `json:"updated_at"`
}

type AppConstraints struct {
	MaxShards null.Int `json:"max_shards,omitzero"`
	MaxGuilds null.Int `json:"max_guilds,omitzero"`
}

func (a AppConstraints) Merge(other AppConstraints) AppConstraints {
	if other.MaxShards.Valid {
		a.MaxShards = other.MaxShards
	}
	if other.MaxGuilds.Valid {
		a.MaxGuilds = other.MaxGuilds
	}
	return a
}

type AppConfig struct {
	Intents  null.Int           `json:"intents,omitzero"`
	Presence *AppPresenceConfig `json:"presence,omitempty"`
}

func (a AppConfig) Merge(other AppConfig) AppConfig {
	if other.Intents.Valid {
		a.Intents = other.Intents
	}
	if other.Presence != nil {
		a.Presence = other.Presence
	}
	return a
}

type AppPresenceConfig struct {
	Status   null.String                `json:"status,omitzero"`
	Activity *AppPresenceActivityConfig `json:"activity,omitempty"`
}

type AppPresenceActivityConfig struct {
	Name  string `json:"name,omitempty"`
	State string `json:"state,omitempty"`
	Type  string `json:"type,omitempty"`
	URL   string `json:"url,omitempty"`
}
