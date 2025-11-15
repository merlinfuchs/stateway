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

func computeChannelPermissions(channel discord.GuildChannel, userID snowflake.ID, roleIDs []snowflake.ID, permissions discord.Permissions) discord.Permissions {
	if overwrite, ok := channel.PermissionOverwrites().Role(channel.GuildID()); ok {
		permissions |= overwrite.Allow
		permissions &= ^overwrite.Deny
	}

	var (
		allow discord.Permissions
		deny  discord.Permissions
	)

	for _, roleID := range roleIDs {
		if roleID == channel.GuildID() {
			continue
		}

		if overwrite, ok := channel.PermissionOverwrites().Role(roleID); ok {
			allow |= overwrite.Allow
			deny |= overwrite.Deny
		}
	}

	if overwrite, ok := channel.PermissionOverwrites().Member(userID); ok {
		allow |= overwrite.Allow
		deny |= overwrite.Deny
	}

	permissions &= ^deny
	permissions |= allow

	return permissions
}
