package model

import (
	"time"

	"github.com/disgoorg/snowflake/v2"
	"gopkg.in/guregu/null.v4"
)

type AppDisabledCode string

const (
	AppDisabledCodeUnknown      AppDisabledCode = "unknown"
	AppDisabledCodeInvalidToken AppDisabledCode = "invalid_token"
)

type App struct {
	ID                  snowflake.ID    `json:"id"`
	GroupID             string          `json:"group_id"`
	DisplayName         string          `json:"display_name"`
	DiscordClientID     snowflake.ID    `json:"discord_client_id"`
	DiscordBotToken     string          `json:"discord_bot_token"`
	DiscordPublicKey    string          `json:"discord_public_key"`
	DiscordClientSecret null.String     `json:"discord_client_secret"`
	Disabled            bool            `json:"disabled"`
	DisabledCode        AppDisabledCode `json:"disabled_code"`
	DisabledMessage     null.String     `json:"disabled_message"`
	CreatedAt           time.Time       `json:"created_at"`
	UpdatedAt           time.Time       `json:"updated_at"`
}

type AppConstraints struct {
	MaxShards null.Int `json:"max_shards,omitzero"`
	MaxGuilds null.Int `json:"max_guilds,omitzero"`
	Intents   null.Int `json:"intents,omitzero"`
}
