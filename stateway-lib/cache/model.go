package cache

import (
	"encoding/json"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

type Channel struct {
	AppID     snowflake.ID    `json:"app_id"`
	GuildID   snowflake.ID    `json:"guild_id"`
	ChannelID snowflake.ID    `json:"channel_id"`
	Data      discord.Channel `json:"data"`
	Tainted   bool            `json:"tainted"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type Emoji struct {
	AppID     snowflake.ID  `json:"app_id"`
	GuildID   snowflake.ID  `json:"guild_id"`
	EmojiID   snowflake.ID  `json:"emoji_id"`
	Data      discord.Emoji `json:"data"`
	Tainted   bool          `json:"tainted"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type Guild struct {
	AppID       snowflake.ID  `json:"app_id"`
	GuildID     snowflake.ID  `json:"guild_id"`
	Data        discord.Guild `json:"data"`
	Unavailable bool          `json:"unavailable"`
	Tainted     bool          `json:"tainted"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

type Role struct {
	AppID     snowflake.ID `json:"app_id"`
	GuildID   snowflake.ID `json:"guild_id"`
	RoleID    snowflake.ID `json:"role_id"`
	Data      discord.Role `json:"data"`
	Tainted   bool         `json:"tainted"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type Sticker struct {
	AppID     snowflake.ID    `json:"app_id"`
	GuildID   snowflake.ID    `json:"guild_id"`
	StickerID snowflake.ID    `json:"sticker_id"`
	Data      discord.Sticker `json:"data"`
	Tainted   bool            `json:"tainted"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

func (c *Channel) UnmarshalJSON(data []byte) error {
	var aux struct {
		AppID     snowflake.ID             `json:"app_id"`
		GuildID   snowflake.ID             `json:"guild_id"`
		ChannelID snowflake.ID             `json:"channel_id"`
		Data      discord.UnmarshalChannel `json:"data"`
		Tainted   bool                     `json:"tainted"`
		CreatedAt time.Time                `json:"created_at"`
		UpdatedAt time.Time                `json:"updated_at"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	c.AppID = aux.AppID
	c.GuildID = aux.GuildID
	c.ChannelID = aux.ChannelID
	c.Data = aux.Data.Channel
	c.Tainted = aux.Tainted
	c.CreatedAt = aux.CreatedAt
	c.UpdatedAt = aux.UpdatedAt
	return nil
}

type CacheOptions struct {
	AppID  snowflake.ID `json:"app_id"`
	Limit  int          `json:"limit"`
	Offset int          `json:"offset"`
}

func ResolveOptions(opts ...CacheOption) CacheOptions {
	options := CacheOptions{}
	for _, opt := range opts {
		opt(&options)
	}
	return options
}

func (o CacheOptions) Destructure() []CacheOption {
	res := []CacheOption{}
	if o.AppID != 0 {
		res = append(res, WithAppID(o.AppID))
	}
	return res
}

type CacheOption func(*CacheOptions)

func WithAppID(appID snowflake.ID) CacheOption {
	return func(o *CacheOptions) {
		o.AppID = appID
	}
}

func WithLimit(limit int) CacheOption {
	return func(o *CacheOptions) {
		o.Limit = limit
	}
}

func WithOffset(offset int) CacheOption {
	return func(o *CacheOptions) {
		o.Offset = offset
	}
}
