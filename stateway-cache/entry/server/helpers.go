package server

import (
	"encoding/json"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/tidwall/sjson"
)

func ensureChannelGuildID(channel discord.GuildChannel, guildID snowflake.ID) discord.GuildChannel {
	return &channelWithGuildID{
		guildID:      guildID,
		GuildChannel: channel,
	}
}

type channelWithGuildID struct {
	guildID snowflake.ID
	discord.GuildChannel
}

func (c *channelWithGuildID) MarshalJSON() ([]byte, error) {
	res, err := json.Marshal(c.GuildChannel)
	if err != nil {
		return nil, err
	}

	res, err = sjson.SetBytes(res, "guild_id", c.guildID.String())
	if err != nil {
		return nil, err
	}

	return res, nil
}
