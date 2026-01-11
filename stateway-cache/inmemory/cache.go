package inmemory

import (
	"context"
	"fmt"
	"time"

	"github.com/disgoorg/snowflake/v2"
	"github.com/hashicorp/go-memdb"
	"github.com/merlinfuchs/stateway/stateway-cache/model"
	"github.com/merlinfuchs/stateway/stateway-cache/store"
)

var cacheSchema = &memdb.DBSchema{
	Tables: map[string]*memdb.TableSchema{
		"guilds": {
			Name: "guilds",
			Indexes: map[string]*memdb.IndexSchema{
				"id": {
					Name:   "id",
					Unique: true,
					Indexer: &memdb.CompoundIndex{Indexes: []memdb.Indexer{
						&memdb.UintFieldIndex{Field: "AppID"},
						&memdb.UintFieldIndex{Field: "GuildID"},
					}},
				},
				"app_id": {
					Name:   "app_id",
					Unique: false,
					Indexer: &memdb.CompoundIndex{Indexes: []memdb.Indexer{
						&memdb.UintFieldIndex{Field: "AppID"},
					}},
				},
			},
		},
		"channels": {
			Name: "channels",
			Indexes: map[string]*memdb.IndexSchema{
				"id": {
					Name:   "id",
					Unique: true,
					Indexer: &memdb.CompoundIndex{Indexes: []memdb.Indexer{
						&memdb.UintFieldIndex{Field: "AppID"},
						&memdb.UintFieldIndex{Field: "GuildID"},
						&memdb.UintFieldIndex{Field: "ChannelID"},
					}},
				},
				"app_id": {
					Name:   "app_id",
					Unique: false,
					Indexer: &memdb.CompoundIndex{Indexes: []memdb.Indexer{
						&memdb.UintFieldIndex{Field: "AppID"},
					}},
				},
				"guild_id": {
					Name:   "guild_id",
					Unique: false,
					Indexer: &memdb.CompoundIndex{Indexes: []memdb.Indexer{
						&memdb.UintFieldIndex{Field: "GuildID"},
					}},
				},
				"channel_id": {
					Name:   "channel_id",
					Unique: false,
					Indexer: &memdb.CompoundIndex{Indexes: []memdb.Indexer{
						&memdb.UintFieldIndex{Field: "AppID"},
						&memdb.UintFieldIndex{Field: "ChannelID"},
					}},
				},
			},
		},
		"roles": {
			Name: "roles",
			Indexes: map[string]*memdb.IndexSchema{
				"id": {
					Name:   "id",
					Unique: true,
					Indexer: &memdb.CompoundIndex{Indexes: []memdb.Indexer{
						&memdb.UintFieldIndex{Field: "AppID"},
						&memdb.UintFieldIndex{Field: "GuildID"},
						&memdb.UintFieldIndex{Field: "RoleID"},
					}},
				},
				"app_id": {
					Name:   "app_id",
					Unique: false,
					Indexer: &memdb.CompoundIndex{Indexes: []memdb.Indexer{
						&memdb.UintFieldIndex{Field: "AppID"},
					}},
				},
				"guild_id": {
					Name:   "guild_id",
					Unique: false,
					Indexer: &memdb.CompoundIndex{Indexes: []memdb.Indexer{
						&memdb.UintFieldIndex{Field: "GuildID"},
					}},
				},
				"role_id": {
					Name:   "role_id",
					Unique: false,
					Indexer: &memdb.CompoundIndex{Indexes: []memdb.Indexer{
						&memdb.UintFieldIndex{Field: "AppID"},
						&memdb.UintFieldIndex{Field: "RoleID"},
					}},
				},
			},
		},
		"emojis": {
			Name: "emojis",
			Indexes: map[string]*memdb.IndexSchema{
				"id": {
					Name:   "id",
					Unique: true,
					Indexer: &memdb.CompoundIndex{Indexes: []memdb.Indexer{
						&memdb.UintFieldIndex{Field: "AppID"},
						&memdb.UintFieldIndex{Field: "GuildID"},
						&memdb.UintFieldIndex{Field: "EmojiID"},
					}},
				},
				"app_id": {
					Name:   "app_id",
					Unique: false,
					Indexer: &memdb.CompoundIndex{Indexes: []memdb.Indexer{
						&memdb.UintFieldIndex{Field: "AppID"},
					}},
				},
				"guild_id": {
					Name:   "guild_id",
					Unique: false,
					Indexer: &memdb.CompoundIndex{Indexes: []memdb.Indexer{
						&memdb.UintFieldIndex{Field: "GuildID"},
					}},
				},
				"emoji_id": {
					Name:   "emoji_id",
					Unique: false,
					Indexer: &memdb.CompoundIndex{Indexes: []memdb.Indexer{
						&memdb.UintFieldIndex{Field: "AppID"},
						&memdb.UintFieldIndex{Field: "EmojiID"},
					}},
				},
			},
		},
		"stickers": {
			Name: "stickers",
			Indexes: map[string]*memdb.IndexSchema{
				"id": {
					Name:   "id",
					Unique: true,
					Indexer: &memdb.CompoundIndex{Indexes: []memdb.Indexer{
						&memdb.UintFieldIndex{Field: "AppID"},
						&memdb.UintFieldIndex{Field: "GuildID"},
						&memdb.UintFieldIndex{Field: "StickerID"},
					}},
				},
				"app_id": {
					Name:   "app_id",
					Unique: false,
					Indexer: &memdb.CompoundIndex{Indexes: []memdb.Indexer{
						&memdb.UintFieldIndex{Field: "AppID"},
					}},
				},
				"guild_id": {
					Name:   "guild_id",
					Unique: false,
					Indexer: &memdb.CompoundIndex{Indexes: []memdb.Indexer{
						&memdb.UintFieldIndex{Field: "GuildID"},
					}},
				},
				"sticker_id": {
					Name:   "sticker_id",
					Unique: false,
					Indexer: &memdb.CompoundIndex{Indexes: []memdb.Indexer{
						&memdb.UintFieldIndex{Field: "AppID"},
						&memdb.UintFieldIndex{Field: "StickerID"},
					}},
				},
			},
		},
	},
}

var _ store.CacheStore = (*InMemoryCacheStore)(nil)

type InMemoryCacheStore struct {
	db *memdb.MemDB
}

func NewInMemoryCacheStore() (*InMemoryCacheStore, error) {
	db, err := memdb.NewMemDB(cacheSchema)
	if err != nil {
		return nil, err
	}
	return &InMemoryCacheStore{db: db}, nil
}

func (s *InMemoryCacheStore) GetGuild(ctx context.Context, appID snowflake.ID, guildID snowflake.ID) (*model.Guild, error) {
	txn := s.db.Txn(false)
	defer txn.Abort()

	guild, err := txn.First("guilds", "id", appID, guildID)
	if err != nil {
		return nil, err
	}

	if guild == nil {
		return nil, store.ErrNotFound
	}

	return guild.(*model.Guild), nil
}

func (s *InMemoryCacheStore) GetGuildOwnerID(ctx context.Context, appID snowflake.ID, guildID snowflake.ID) (snowflake.ID, error) {
	txn := s.db.Txn(false)
	defer txn.Abort()

	guild, err := txn.First("guilds", "id", appID, guildID)
	if err != nil {
		return 0, err
	}

	if guild == nil {
		return 0, store.ErrNotFound
	}

	return guild.(*model.Guild).Data.OwnerID, nil
}

func (s *InMemoryCacheStore) GetGuilds(ctx context.Context, appID snowflake.ID, limit int, offset int) ([]*model.Guild, error) {
	txn := s.db.Txn(false)
	defer txn.Abort()

	iter, err := txn.Get("guilds", "app_id", appID)
	if err != nil {
		return nil, err
	}

	guilds := make([]*model.Guild, 0, limit)
	for guild := iter.Next(); guild != nil; guild = iter.Next() {
		if offset > 0 {
			offset--
			continue
		}
		if limit > 0 && len(guilds) >= limit {
			break
		}
		guilds = append(guilds, guild.(*model.Guild))
	}

	return guilds, nil
}

func (s *InMemoryCacheStore) CheckGuildExist(ctx context.Context, appID snowflake.ID, guildID snowflake.ID) (bool, error) {
	txn := s.db.Txn(false)
	defer txn.Abort()

	guild, err := txn.First("guilds", "id", appID, guildID)
	if err != nil {
		return false, err
	}

	return guild != nil, nil
}

func (s *InMemoryCacheStore) UpsertGuilds(ctx context.Context, guilds ...store.UpsertGuildParams) error {
	txn := s.db.Txn(true)
	defer txn.Abort()

	for _, guild := range guilds {
		createdAt := guild.CreatedAt
		if createdAt.IsZero() {
			createdAt = time.Now().UTC()
		}
		updatedAt := guild.UpdatedAt
		if updatedAt.IsZero() {
			updatedAt = time.Now().UTC()
		}

		err := txn.Insert("guilds", &model.Guild{
			AppID:     guild.AppID,
			GuildID:   guild.GuildID,
			Data:      guild.Data,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
		if err != nil {
			return err
		}
	}

	txn.Commit()
	return nil
}

func (s *InMemoryCacheStore) MarkGuildUnavailable(ctx context.Context, appID snowflake.ID, guildID snowflake.ID) error {
	txn := s.db.Txn(true)
	defer txn.Abort()

	guild, err := txn.First("guilds", "id", appID, guildID)
	if err != nil {
		return err
	}

	if guild == nil {
		return nil
	}

	g := *guild.(*model.Guild)
	g.Unavailable = true
	g.UpdatedAt = time.Now().UTC()

	err = txn.Insert("guilds", &g)
	if err != nil {
		return err
	}

	txn.Commit()
	return nil
}

func (s *InMemoryCacheStore) DeleteGuild(ctx context.Context, appID snowflake.ID, guildID snowflake.ID) error {
	txn := s.db.Txn(true)
	defer txn.Abort()

	_, err := txn.DeleteAll("guilds", "id", appID, guildID)
	if err != nil {
		return err
	}

	txn.Commit()
	return nil
}

func (s *InMemoryCacheStore) SearchGuilds(ctx context.Context, params store.SearchGuildsParams) ([]*model.Guild, error) {
	return nil, fmt.Errorf("not implemented")
}

// CacheRoleStore methods

func (s *InMemoryCacheStore) GetGuildRole(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, roleID snowflake.ID) (*model.Role, error) {
	txn := s.db.Txn(false)
	defer txn.Abort()

	role, err := txn.First("roles", "id", appID, guildID, roleID)
	if err != nil {
		return nil, err
	}

	if role == nil {
		return nil, store.ErrNotFound
	}

	return role.(*model.Role), nil
}

func (s *InMemoryCacheStore) GetRole(ctx context.Context, appID snowflake.ID, roleID snowflake.ID) (*model.Role, error) {
	txn := s.db.Txn(false)
	defer txn.Abort()

	role, err := txn.First("roles", "role_id", appID, roleID)
	if err != nil {
		return nil, err
	}

	if role == nil {
		return nil, store.ErrNotFound
	}

	return role.(*model.Role), nil
}

func (s *InMemoryCacheStore) GetGuildRoles(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, limit int, offset int) ([]*model.Role, error) {
	txn := s.db.Txn(false)
	defer txn.Abort()

	iter, err := txn.Get("roles", "guild_id", guildID)
	if err != nil {
		return nil, err
	}

	roles := make([]*model.Role, 0, limit)
	for role := iter.Next(); role != nil; role = iter.Next() {
		r := role.(*model.Role)
		if r.AppID != appID {
			continue
		}
		if offset > 0 {
			offset--
			continue
		}
		if limit > 0 && len(roles) >= limit {
			break
		}
		roles = append(roles, r)
	}

	return roles, nil
}

func (s *InMemoryCacheStore) GetGuildRolesByIDs(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, roleIDs []snowflake.ID) ([]*model.Role, error) {
	txn := s.db.Txn(false)
	defer txn.Abort()

	roles := make([]*model.Role, 0, len(roleIDs))
	for _, roleID := range roleIDs {
		role, err := txn.First("roles", "id", appID, guildID, roleID)
		if err != nil {
			return nil, err
		}
		if role != nil {
			roles = append(roles, role.(*model.Role))
		}
	}

	return roles, nil
}

func (s *InMemoryCacheStore) GetRoles(ctx context.Context, appID snowflake.ID, limit int, offset int) ([]*model.Role, error) {
	txn := s.db.Txn(false)
	defer txn.Abort()

	iter, err := txn.Get("roles", "app_id", appID)
	if err != nil {
		return nil, err
	}

	roles := make([]*model.Role, 0, limit)
	for role := iter.Next(); role != nil; role = iter.Next() {
		if offset > 0 {
			offset--
			continue
		}
		if limit > 0 && len(roles) >= limit {
			break
		}
		roles = append(roles, role.(*model.Role))
	}

	return roles, nil
}

func (s *InMemoryCacheStore) SearchGuildRoles(ctx context.Context, params store.SearchGuildRolesParams) ([]*model.Role, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *InMemoryCacheStore) SearchRoles(ctx context.Context, params store.SearchRolesParams) ([]*model.Role, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *InMemoryCacheStore) CountGuildRoles(ctx context.Context, appID snowflake.ID, guildID snowflake.ID) (int, error) {
	txn := s.db.Txn(false)
	defer txn.Abort()

	iter, err := txn.Get("roles", "guild_id", guildID)
	if err != nil {
		return 0, err
	}

	var count int
	for role := iter.Next(); role != nil; role = iter.Next() {
		if role.(*model.Role).AppID == appID {
			count++
		}
	}

	return count, nil
}

func (s *InMemoryCacheStore) CountRoles(ctx context.Context, appID snowflake.ID) (int, error) {
	txn := s.db.Txn(false)
	defer txn.Abort()

	iter, err := txn.Get("roles", "app_id", appID)
	if err != nil {
		return 0, err
	}

	var count int
	for iter.Next() != nil {
		count++
	}

	return count, nil
}

func (s *InMemoryCacheStore) UpsertRoles(ctx context.Context, roles ...store.UpsertRoleParams) error {
	txn := s.db.Txn(true)
	defer txn.Abort()

	for _, role := range roles {
		createdAt := role.CreatedAt
		if createdAt.IsZero() {
			createdAt = time.Now().UTC()
		}
		updatedAt := role.UpdatedAt
		if updatedAt.IsZero() {
			updatedAt = time.Now().UTC()
		}

		err := txn.Insert("roles", &model.Role{
			AppID:     role.AppID,
			GuildID:   role.GuildID,
			RoleID:    role.RoleID,
			Data:      role.Data,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
		if err != nil {
			return err
		}
	}

	txn.Commit()
	return nil
}

func (s *InMemoryCacheStore) DeleteRole(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, roleID snowflake.ID) error {
	txn := s.db.Txn(true)
	defer txn.Abort()

	_, err := txn.DeleteAll("roles", "id", appID, guildID, roleID)
	if err != nil {
		return err
	}

	txn.Commit()
	return nil
}

// CacheChannelStore methods

func (s *InMemoryCacheStore) GetGuildChannel(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, channelID snowflake.ID) (*model.Channel, error) {
	txn := s.db.Txn(false)
	defer txn.Abort()

	channel, err := txn.First("channels", "id", appID, guildID, channelID)
	if err != nil {
		return nil, err
	}

	if channel == nil {
		return nil, store.ErrNotFound
	}

	return channel.(*model.Channel), nil
}

func (s *InMemoryCacheStore) GetChannel(ctx context.Context, appID snowflake.ID, channelID snowflake.ID) (*model.Channel, error) {
	txn := s.db.Txn(false)
	defer txn.Abort()

	channel, err := txn.First("channels", "channel_id", appID, channelID)
	if err != nil {
		return nil, err
	}

	if channel == nil {
		return nil, store.ErrNotFound
	}

	return channel.(*model.Channel), nil
}

func (s *InMemoryCacheStore) GetGuildChannels(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, limit int, offset int) ([]*model.Channel, error) {
	txn := s.db.Txn(false)
	defer txn.Abort()

	iter, err := txn.Get("channels", "guild_id", guildID)
	if err != nil {
		return nil, err
	}

	channels := make([]*model.Channel, 0, limit)
	for channel := iter.Next(); channel != nil; channel = iter.Next() {
		c := channel.(*model.Channel)
		if c.AppID != appID {
			continue
		}
		if offset > 0 {
			offset--
			continue
		}
		if limit > 0 && len(channels) >= limit {
			break
		}
		channels = append(channels, c)
	}

	return channels, nil
}

func (s *InMemoryCacheStore) GetChannels(ctx context.Context, appID snowflake.ID, limit int, offset int) ([]*model.Channel, error) {
	txn := s.db.Txn(false)
	defer txn.Abort()

	iter, err := txn.Get("channels", "app_id", appID)
	if err != nil {
		return nil, err
	}

	channels := make([]*model.Channel, 0, limit)
	for channel := iter.Next(); channel != nil; channel = iter.Next() {
		if offset > 0 {
			offset--
			continue
		}
		if limit > 0 && len(channels) >= limit {
			break
		}
		channels = append(channels, channel.(*model.Channel))
	}

	return channels, nil
}

func (s *InMemoryCacheStore) GetChannelsByType(ctx context.Context, appID snowflake.ID, types []int, limit int, offset int) ([]*model.Channel, error) {
	txn := s.db.Txn(false)
	defer txn.Abort()

	iter, err := txn.Get("channels", "app_id", appID)
	if err != nil {
		return nil, err
	}

	channels := make([]*model.Channel, 0, limit)
	for channel := iter.Next(); channel != nil; channel = iter.Next() {
		c := channel.(*model.Channel)
		channelType := int(c.Data.Type())
		found := false
		for _, t := range types {
			if t == channelType {
				found = true
				break
			}
		}
		if !found {
			continue
		}
		if offset > 0 {
			offset--
			continue
		}
		if limit > 0 && len(channels) >= limit {
			break
		}
		channels = append(channels, c)
	}

	return channels, nil
}

func (s *InMemoryCacheStore) SearchGuildChannels(ctx context.Context, params store.SearchGuildChannelsParams) ([]*model.Channel, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *InMemoryCacheStore) SearchChannels(ctx context.Context, params store.SearchChannelsParams) ([]*model.Channel, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *InMemoryCacheStore) CountGuildChannels(ctx context.Context, appID snowflake.ID, guildID snowflake.ID) (int, error) {
	txn := s.db.Txn(false)
	defer txn.Abort()

	iter, err := txn.Get("channels", "guild_id", guildID)
	if err != nil {
		return 0, err
	}

	var count int
	for channel := iter.Next(); channel != nil; channel = iter.Next() {
		if channel.(*model.Channel).AppID == appID {
			count++
		}
	}

	return count, nil
}

func (s *InMemoryCacheStore) CountChannels(ctx context.Context, appID snowflake.ID) (int, error) {
	txn := s.db.Txn(false)
	defer txn.Abort()

	iter, err := txn.Get("channels", "app_id", appID)
	if err != nil {
		return 0, err
	}

	var count int
	for iter.Next() != nil {
		count++
	}

	return count, nil
}

func (s *InMemoryCacheStore) UpsertChannels(ctx context.Context, channels ...store.UpsertChannelParams) error {
	txn := s.db.Txn(true)
	defer txn.Abort()

	for _, channel := range channels {
		createdAt := channel.CreatedAt
		if createdAt.IsZero() {
			createdAt = time.Now().UTC()
		}
		updatedAt := channel.UpdatedAt
		if updatedAt.IsZero() {
			updatedAt = time.Now().UTC()
		}

		err := txn.Insert("channels", &model.Channel{
			AppID:     channel.AppID,
			GuildID:   channel.GuildID,
			ChannelID: channel.ChannelID,
			Data:      channel.Data,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
		if err != nil {
			return err
		}
	}

	txn.Commit()
	return nil
}

func (s *InMemoryCacheStore) DeleteChannel(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, channelID snowflake.ID) error {
	txn := s.db.Txn(true)
	defer txn.Abort()

	_, err := txn.DeleteAll("channels", "id", appID, guildID, channelID)
	if err != nil {
		return err
	}

	txn.Commit()
	return nil
}

// CacheEmojiStore methods

func (s *InMemoryCacheStore) GetGuildEmoji(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, emojiID snowflake.ID) (*model.Emoji, error) {
	txn := s.db.Txn(false)
	defer txn.Abort()

	emoji, err := txn.First("emojis", "id", appID, guildID, emojiID)
	if err != nil {
		return nil, err
	}

	if emoji == nil {
		return nil, store.ErrNotFound
	}

	return emoji.(*model.Emoji), nil
}

func (s *InMemoryCacheStore) GetEmoji(ctx context.Context, appID snowflake.ID, emojiID snowflake.ID) (*model.Emoji, error) {
	txn := s.db.Txn(false)
	defer txn.Abort()

	emoji, err := txn.First("emojis", "emoji_id", appID, emojiID)
	if err != nil {
		return nil, err
	}

	if emoji == nil {
		return nil, store.ErrNotFound
	}

	return emoji.(*model.Emoji), nil
}

func (s *InMemoryCacheStore) GetGuildEmojis(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, limit int, offset int) ([]*model.Emoji, error) {
	txn := s.db.Txn(false)
	defer txn.Abort()

	iter, err := txn.Get("emojis", "guild_id", guildID)
	if err != nil {
		return nil, err
	}

	emojis := make([]*model.Emoji, 0, limit)
	for emoji := iter.Next(); emoji != nil; emoji = iter.Next() {
		e := emoji.(*model.Emoji)
		if e.AppID != appID {
			continue
		}
		if offset > 0 {
			offset--
			continue
		}
		if limit > 0 && len(emojis) >= limit {
			break
		}
		emojis = append(emojis, e)
	}

	return emojis, nil
}

func (s *InMemoryCacheStore) GetEmojis(ctx context.Context, appID snowflake.ID, limit int, offset int) ([]*model.Emoji, error) {
	txn := s.db.Txn(false)
	defer txn.Abort()

	iter, err := txn.Get("emojis", "app_id", appID)
	if err != nil {
		return nil, err
	}

	emojis := make([]*model.Emoji, 0, limit)
	for emoji := iter.Next(); emoji != nil; emoji = iter.Next() {
		if offset > 0 {
			offset--
			continue
		}
		if limit > 0 && len(emojis) >= limit {
			break
		}
		emojis = append(emojis, emoji.(*model.Emoji))
	}

	return emojis, nil
}

func (s *InMemoryCacheStore) SearchGuildEmojis(ctx context.Context, params store.SearchGuildEmojisParams) ([]*model.Emoji, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *InMemoryCacheStore) SearchEmojis(ctx context.Context, params store.SearchEmojisParams) ([]*model.Emoji, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *InMemoryCacheStore) CountGuildEmojis(ctx context.Context, appID snowflake.ID, guildID snowflake.ID) (int, error) {
	txn := s.db.Txn(false)
	defer txn.Abort()

	iter, err := txn.Get("emojis", "guild_id", guildID)
	if err != nil {
		return 0, err
	}

	var count int
	for emoji := iter.Next(); emoji != nil; emoji = iter.Next() {
		if emoji.(*model.Emoji).AppID == appID {
			count++
		}
	}

	return count, nil
}

func (s *InMemoryCacheStore) CountEmojis(ctx context.Context, appID snowflake.ID) (int, error) {
	txn := s.db.Txn(false)
	defer txn.Abort()

	iter, err := txn.Get("emojis", "app_id", appID)
	if err != nil {
		return 0, err
	}

	var count int
	for iter.Next() != nil {
		count++
	}

	return count, nil
}

func (s *InMemoryCacheStore) UpsertEmojis(ctx context.Context, emojis ...store.UpsertEmojiParams) error {
	txn := s.db.Txn(true)
	defer txn.Abort()

	for _, emoji := range emojis {
		createdAt := emoji.CreatedAt
		if createdAt.IsZero() {
			createdAt = time.Now().UTC()
		}
		updatedAt := emoji.UpdatedAt
		if updatedAt.IsZero() {
			updatedAt = time.Now().UTC()
		}

		err := txn.Insert("emojis", &model.Emoji{
			AppID:     emoji.AppID,
			GuildID:   emoji.GuildID,
			EmojiID:   emoji.EmojiID,
			Data:      emoji.Data,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
		if err != nil {
			return err
		}
	}

	txn.Commit()
	return nil
}

func (s *InMemoryCacheStore) DeleteEmoji(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, emojiID snowflake.ID) error {
	txn := s.db.Txn(true)
	defer txn.Abort()

	_, err := txn.DeleteAll("emojis", "id", appID, guildID, emojiID)
	if err != nil {
		return err
	}

	txn.Commit()
	return nil
}

// CacheStickerStore methods

func (s *InMemoryCacheStore) GetSticker(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, stickerID snowflake.ID) (*model.Sticker, error) {
	txn := s.db.Txn(false)
	defer txn.Abort()

	sticker, err := txn.First("stickers", "id", appID, guildID, stickerID)
	if err != nil {
		return nil, err
	}

	if sticker == nil {
		return nil, store.ErrNotFound
	}

	return sticker.(*model.Sticker), nil
}

func (s *InMemoryCacheStore) GetGuildSticker(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, stickerID snowflake.ID) (*model.Sticker, error) {
	return s.GetSticker(ctx, appID, guildID, stickerID)
}

func (s *InMemoryCacheStore) GetStickers(ctx context.Context, appID snowflake.ID, limit int, offset int) ([]*model.Sticker, error) {
	txn := s.db.Txn(false)
	defer txn.Abort()

	iter, err := txn.Get("stickers", "app_id", appID)
	if err != nil {
		return nil, err
	}

	stickers := make([]*model.Sticker, 0, limit)
	for sticker := iter.Next(); sticker != nil; sticker = iter.Next() {
		if offset > 0 {
			offset--
			continue
		}
		if limit > 0 && len(stickers) >= limit {
			break
		}
		stickers = append(stickers, sticker.(*model.Sticker))
	}

	return stickers, nil
}

func (s *InMemoryCacheStore) GetGuildStickers(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, limit int, offset int) ([]*model.Sticker, error) {
	txn := s.db.Txn(false)
	defer txn.Abort()

	iter, err := txn.Get("stickers", "guild_id", guildID)
	if err != nil {
		return nil, err
	}

	stickers := make([]*model.Sticker, 0, limit)
	for sticker := iter.Next(); sticker != nil; sticker = iter.Next() {
		s := sticker.(*model.Sticker)
		if s.AppID != appID {
			continue
		}
		if offset > 0 {
			offset--
			continue
		}
		if limit > 0 && len(stickers) >= limit {
			break
		}
		stickers = append(stickers, s)
	}

	return stickers, nil
}

func (s *InMemoryCacheStore) SearchStickers(ctx context.Context, params store.SearchStickersParams) ([]*model.Sticker, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *InMemoryCacheStore) SearchGuildStickers(ctx context.Context, params store.SearchGuildStickersParams) ([]*model.Sticker, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *InMemoryCacheStore) CountGuildStickers(ctx context.Context, appID snowflake.ID, guildID snowflake.ID) (int, error) {
	txn := s.db.Txn(false)
	defer txn.Abort()

	iter, err := txn.Get("stickers", "guild_id", guildID)
	if err != nil {
		return 0, err
	}

	var count int
	for sticker := iter.Next(); sticker != nil; sticker = iter.Next() {
		if sticker.(*model.Sticker).AppID == appID {
			count++
		}
	}

	return count, nil
}

func (s *InMemoryCacheStore) CountStickers(ctx context.Context, appID snowflake.ID) (int, error) {
	txn := s.db.Txn(false)
	defer txn.Abort()

	iter, err := txn.Get("stickers", "app_id", appID)
	if err != nil {
		return 0, err
	}

	var count int
	for iter.Next() != nil {
		count++
	}

	return count, nil
}

func (s *InMemoryCacheStore) UpsertStickers(ctx context.Context, stickers ...store.UpsertStickerParams) error {
	txn := s.db.Txn(true)
	defer txn.Abort()

	for _, sticker := range stickers {
		createdAt := sticker.CreatedAt
		if createdAt.IsZero() {
			createdAt = time.Now().UTC()
		}
		updatedAt := sticker.UpdatedAt
		if updatedAt.IsZero() {
			updatedAt = time.Now().UTC()
		}

		err := txn.Insert("stickers", &model.Sticker{
			AppID:     sticker.AppID,
			GuildID:   sticker.GuildID,
			StickerID: sticker.StickerID,
			Data:      sticker.Data,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
		if err != nil {
			return err
		}
	}

	txn.Commit()
	return nil
}

func (s *InMemoryCacheStore) DeleteSticker(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, stickerID snowflake.ID) error {
	txn := s.db.Txn(true)
	defer txn.Abort()

	_, err := txn.DeleteAll("stickers", "id", appID, guildID, stickerID)
	if err != nil {
		return err
	}

	txn.Commit()
	return nil
}

// CacheStore methods

func (s *InMemoryCacheStore) MarkShardEntitiesTainted(ctx context.Context, params store.MarkShardEntitiesTaintedParams) error {
	// Not implemented for in-memory store
	return nil
}

func (s *InMemoryCacheStore) MassUpsertEntities(ctx context.Context, params store.MassUpsertEntitiesParams) error {
	txn := s.db.Txn(true)
	defer txn.Abort()

	for _, guild := range params.Guilds {
		createdAt := guild.CreatedAt
		if createdAt.IsZero() {
			createdAt = time.Now().UTC()
		}
		updatedAt := guild.UpdatedAt
		if updatedAt.IsZero() {
			updatedAt = time.Now().UTC()
		}

		err := txn.Insert("guilds", &model.Guild{
			AppID:     guild.AppID,
			GuildID:   guild.GuildID,
			Data:      guild.Data,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
		if err != nil {
			return err
		}
	}

	for _, role := range params.Roles {
		createdAt := role.CreatedAt
		if createdAt.IsZero() {
			createdAt = time.Now().UTC()
		}
		updatedAt := role.UpdatedAt
		if updatedAt.IsZero() {
			updatedAt = time.Now().UTC()
		}

		err := txn.Insert("roles", &model.Role{
			AppID:     role.AppID,
			GuildID:   role.GuildID,
			RoleID:    role.RoleID,
			Data:      role.Data,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
		if err != nil {
			return err
		}
	}

	for _, channel := range params.Channels {
		createdAt := channel.CreatedAt
		if createdAt.IsZero() {
			createdAt = time.Now().UTC()
		}
		updatedAt := channel.UpdatedAt
		if updatedAt.IsZero() {
			updatedAt = time.Now().UTC()
		}

		err := txn.Insert("channels", &model.Channel{
			AppID:     channel.AppID,
			GuildID:   channel.GuildID,
			ChannelID: channel.ChannelID,
			Data:      channel.Data,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
		if err != nil {
			return err
		}
	}

	for _, emoji := range params.Emojis {
		createdAt := emoji.CreatedAt
		if createdAt.IsZero() {
			createdAt = time.Now().UTC()
		}
		updatedAt := emoji.UpdatedAt
		if updatedAt.IsZero() {
			updatedAt = time.Now().UTC()
		}

		err := txn.Insert("emojis", &model.Emoji{
			AppID:     emoji.AppID,
			GuildID:   emoji.GuildID,
			EmojiID:   emoji.EmojiID,
			Data:      emoji.Data,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
		if err != nil {
			return err
		}
	}

	for _, sticker := range params.Stickers {
		createdAt := sticker.CreatedAt
		if createdAt.IsZero() {
			createdAt = time.Now().UTC()
		}
		updatedAt := sticker.UpdatedAt
		if updatedAt.IsZero() {
			updatedAt = time.Now().UTC()
		}

		err := txn.Insert("stickers", &model.Sticker{
			AppID:     sticker.AppID,
			GuildID:   sticker.GuildID,
			StickerID: sticker.StickerID,
			Data:      sticker.Data,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
		if err != nil {
			return err
		}
	}

	txn.Commit()
	return nil
}
