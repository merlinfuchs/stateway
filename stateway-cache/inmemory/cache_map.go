package inmemory

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/disgoorg/snowflake/v2"
	"github.com/merlinfuchs/stateway/stateway-cache/model"
	"github.com/merlinfuchs/stateway/stateway-cache/store"
)

var _ store.CacheStore = (*MapCacheStore)(nil)

// MapCacheStore is a thread-safe in-memory cache store implementation using Go maps.
// It uses RWMutex for fine-grained locking and maintains multiple indexes for efficient queries.
type MapCacheStore struct {
	// Guilds: primary index by appID -> guildID
	guildsMu sync.RWMutex
	guilds   map[snowflake.ID]map[snowflake.ID]*model.Guild

	// Channels: primary index by appID -> channelID, guild index by appID -> guildID -> channelID
	channelsMu      sync.RWMutex
	channels        map[snowflake.ID]map[snowflake.ID]*model.Channel
	channelsByGuild map[snowflake.ID]map[snowflake.ID]map[snowflake.ID]*model.Channel

	// Roles: primary index by appID -> roleID, guild index by appID -> guildID -> roleID
	rolesMu      sync.RWMutex
	roles        map[snowflake.ID]map[snowflake.ID]*model.Role
	rolesByGuild map[snowflake.ID]map[snowflake.ID]map[snowflake.ID]*model.Role

	// Emojis: primary index by appID -> emojiID, guild index by appID -> guildID -> emojiID
	emojisMu      sync.RWMutex
	emojis        map[snowflake.ID]map[snowflake.ID]*model.Emoji
	emojisByGuild map[snowflake.ID]map[snowflake.ID]map[snowflake.ID]*model.Emoji

	// Stickers: primary index by appID -> stickerID, guild index by appID -> guildID -> stickerID
	stickersMu      sync.RWMutex
	stickers        map[snowflake.ID]map[snowflake.ID]*model.Sticker
	stickersByGuild map[snowflake.ID]map[snowflake.ID]map[snowflake.ID]*model.Sticker
}

func NewMapCacheStore() *MapCacheStore {
	return &MapCacheStore{
		guilds:          make(map[snowflake.ID]map[snowflake.ID]*model.Guild),
		channels:        make(map[snowflake.ID]map[snowflake.ID]*model.Channel),
		channelsByGuild: make(map[snowflake.ID]map[snowflake.ID]map[snowflake.ID]*model.Channel),
		roles:           make(map[snowflake.ID]map[snowflake.ID]*model.Role),
		rolesByGuild:    make(map[snowflake.ID]map[snowflake.ID]map[snowflake.ID]*model.Role),
		emojis:          make(map[snowflake.ID]map[snowflake.ID]*model.Emoji),
		emojisByGuild:   make(map[snowflake.ID]map[snowflake.ID]map[snowflake.ID]*model.Emoji),
		stickers:        make(map[snowflake.ID]map[snowflake.ID]*model.Sticker),
		stickersByGuild: make(map[snowflake.ID]map[snowflake.ID]map[snowflake.ID]*model.Sticker),
	}
}

// CacheGuildStore methods

func (s *MapCacheStore) GetGuild(ctx context.Context, appID snowflake.ID, guildID snowflake.ID) (*model.Guild, error) {
	s.guildsMu.RLock()
	defer s.guildsMu.RUnlock()

	appGuilds, ok := s.guilds[appID]
	if !ok {
		return nil, store.ErrNotFound
	}

	guild, ok := appGuilds[guildID]
	if !ok {
		return nil, store.ErrNotFound
	}

	return guild, nil
}

func (s *MapCacheStore) GetGuildOwnerID(ctx context.Context, appID snowflake.ID, guildID snowflake.ID) (snowflake.ID, error) {
	guild, err := s.GetGuild(ctx, appID, guildID)
	if err != nil {
		return 0, err
	}
	return guild.Data.OwnerID, nil
}

func (s *MapCacheStore) GetGuilds(ctx context.Context, appID snowflake.ID, limit int, offset int) ([]*model.Guild, error) {
	s.guildsMu.RLock()
	defer s.guildsMu.RUnlock()

	appGuilds, ok := s.guilds[appID]
	if !ok {
		return []*model.Guild{}, nil
	}

	guilds := make([]*model.Guild, 0)
	if limit > 0 {
		guilds = make([]*model.Guild, 0, limit)
	}

	currentOffset := 0
	for _, guild := range appGuilds {
		if currentOffset < offset {
			currentOffset++
			continue
		}

		guilds = append(guilds, guild)

		if limit > 0 && len(guilds) >= limit {
			break
		}
	}

	return guilds, nil
}

func (s *MapCacheStore) CheckGuildExist(ctx context.Context, appID snowflake.ID, guildID snowflake.ID) (bool, error) {
	s.guildsMu.RLock()
	defer s.guildsMu.RUnlock()

	appGuilds, ok := s.guilds[appID]
	if !ok {
		return false, nil
	}

	_, ok = appGuilds[guildID]
	return ok, nil
}

func (s *MapCacheStore) UpsertGuilds(ctx context.Context, guilds ...store.UpsertGuildParams) error {
	s.guildsMu.Lock()
	defer s.guildsMu.Unlock()

	for _, guild := range guilds {
		if s.guilds[guild.AppID] == nil {
			s.guilds[guild.AppID] = make(map[snowflake.ID]*model.Guild)
		}

		createdAt := guild.CreatedAt
		if createdAt.IsZero() {
			createdAt = time.Now().UTC()
		}
		updatedAt := guild.UpdatedAt
		if updatedAt.IsZero() {
			updatedAt = time.Now().UTC()
		}

		s.guilds[guild.AppID][guild.GuildID] = &model.Guild{
			AppID:     guild.AppID,
			GuildID:   guild.GuildID,
			Data:      guild.Data,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}
	}

	return nil
}

func (s *MapCacheStore) MarkGuildUnavailable(ctx context.Context, appID snowflake.ID, guildID snowflake.ID) error {
	s.guildsMu.Lock()
	defer s.guildsMu.Unlock()

	appGuilds, ok := s.guilds[appID]
	if !ok {
		return nil
	}

	guild, ok := appGuilds[guildID]
	if !ok {
		return nil
	}

	guild.Unavailable = true
	guild.UpdatedAt = time.Now().UTC()

	return nil
}

func (s *MapCacheStore) DeleteGuild(ctx context.Context, appID snowflake.ID, guildID snowflake.ID) error {
	s.guildsMu.Lock()
	defer s.guildsMu.Unlock()

	appGuilds, ok := s.guilds[appID]
	if !ok {
		return nil
	}

	_, ok = appGuilds[guildID]
	if !ok {
		return nil
	}

	delete(appGuilds, guildID)
	if len(appGuilds) == 0 {
		delete(s.guilds, appID)
	}

	return nil
}

func (s *MapCacheStore) SearchGuilds(ctx context.Context, params store.SearchGuildsParams) ([]*model.Guild, error) {
	return nil, fmt.Errorf("not implemented")
}

// CacheRoleStore methods

func (s *MapCacheStore) GetGuildRole(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, roleID snowflake.ID) (*model.Role, error) {
	s.rolesMu.RLock()
	defer s.rolesMu.RUnlock()

	guildRoles, ok := s.rolesByGuild[appID]
	if !ok {
		return nil, store.ErrNotFound
	}

	roles, ok := guildRoles[guildID]
	if !ok {
		return nil, store.ErrNotFound
	}

	role, ok := roles[roleID]
	if !ok {
		return nil, store.ErrNotFound
	}

	return role, nil
}

func (s *MapCacheStore) GetRole(ctx context.Context, appID snowflake.ID, roleID snowflake.ID) (*model.Role, error) {
	s.rolesMu.RLock()
	defer s.rolesMu.RUnlock()

	appRoles, ok := s.roles[appID]
	if !ok {
		return nil, store.ErrNotFound
	}

	role, ok := appRoles[roleID]
	if !ok {
		return nil, store.ErrNotFound
	}

	return role, nil
}

func (s *MapCacheStore) GetGuildRoles(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, limit int, offset int) ([]*model.Role, error) {
	s.rolesMu.RLock()
	defer s.rolesMu.RUnlock()

	guildRoles, ok := s.rolesByGuild[appID]
	if !ok {
		return []*model.Role{}, nil
	}

	roles, ok := guildRoles[guildID]
	if !ok {
		return []*model.Role{}, nil
	}

	result := make([]*model.Role, 0)
	if limit > 0 {
		result = make([]*model.Role, 0, limit)
	}

	currentOffset := 0
	for _, role := range roles {
		if currentOffset < offset {
			currentOffset++
			continue
		}

		result = append(result, role)

		if limit > 0 && len(result) >= limit {
			break
		}
	}

	return result, nil
}

func (s *MapCacheStore) GetGuildRolesByIDs(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, roleIDs []snowflake.ID) ([]*model.Role, error) {
	s.rolesMu.RLock()
	defer s.rolesMu.RUnlock()

	guildRoles, ok := s.rolesByGuild[appID]
	if !ok {
		return []*model.Role{}, nil
	}

	roles, ok := guildRoles[guildID]
	if !ok {
		return []*model.Role{}, nil
	}

	result := make([]*model.Role, 0, len(roleIDs))
	for _, roleID := range roleIDs {
		if role, ok := roles[roleID]; ok {
			result = append(result, role)
		}
	}

	return result, nil
}

func (s *MapCacheStore) GetRoles(ctx context.Context, appID snowflake.ID, limit int, offset int) ([]*model.Role, error) {
	s.rolesMu.RLock()
	defer s.rolesMu.RUnlock()

	appRoles, ok := s.roles[appID]
	if !ok {
		return []*model.Role{}, nil
	}

	roles := make([]*model.Role, 0)
	if limit > 0 {
		roles = make([]*model.Role, 0, limit)
	}

	currentOffset := 0
	for _, role := range appRoles {
		if currentOffset < offset {
			currentOffset++
			continue
		}

		roles = append(roles, role)

		if limit > 0 && len(roles) >= limit {
			break
		}
	}

	return roles, nil
}

func (s *MapCacheStore) SearchGuildRoles(ctx context.Context, params store.SearchGuildRolesParams) ([]*model.Role, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *MapCacheStore) SearchRoles(ctx context.Context, params store.SearchRolesParams) ([]*model.Role, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *MapCacheStore) CountGuildRoles(ctx context.Context, appID snowflake.ID, guildID snowflake.ID) (int, error) {
	s.rolesMu.RLock()
	defer s.rolesMu.RUnlock()

	guildRoles, ok := s.rolesByGuild[appID]
	if !ok {
		return 0, nil
	}

	roles, ok := guildRoles[guildID]
	if !ok {
		return 0, nil
	}

	return len(roles), nil
}

func (s *MapCacheStore) CountRoles(ctx context.Context, appID snowflake.ID) (int, error) {
	s.rolesMu.RLock()
	defer s.rolesMu.RUnlock()

	appRoles, ok := s.roles[appID]
	if !ok {
		return 0, nil
	}

	return len(appRoles), nil
}

func (s *MapCacheStore) UpsertRoles(ctx context.Context, roles ...store.UpsertRoleParams) error {
	s.rolesMu.Lock()
	defer s.rolesMu.Unlock()

	for _, role := range roles {
		if s.roles[role.AppID] == nil {
			s.roles[role.AppID] = make(map[snowflake.ID]*model.Role)
		}
		if s.rolesByGuild[role.AppID] == nil {
			s.rolesByGuild[role.AppID] = make(map[snowflake.ID]map[snowflake.ID]*model.Role)
		}
		if s.rolesByGuild[role.AppID][role.GuildID] == nil {
			s.rolesByGuild[role.AppID][role.GuildID] = make(map[snowflake.ID]*model.Role)
		}

		createdAt := role.CreatedAt
		if createdAt.IsZero() {
			createdAt = time.Now().UTC()
		}
		updatedAt := role.UpdatedAt
		if updatedAt.IsZero() {
			updatedAt = time.Now().UTC()
		}

		r := &model.Role{
			AppID:     role.AppID,
			GuildID:   role.GuildID,
			RoleID:    role.RoleID,
			Data:      role.Data,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}

		s.roles[role.AppID][role.RoleID] = r
		s.rolesByGuild[role.AppID][role.GuildID][role.RoleID] = r
	}

	return nil
}

func (s *MapCacheStore) DeleteRole(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, roleID snowflake.ID) error {
	s.rolesMu.Lock()
	defer s.rolesMu.Unlock()

	appRoles, ok := s.roles[appID]
	if !ok {
		return nil
	}

	_, ok = appRoles[roleID]
	if !ok {
		return nil
	}

	delete(appRoles, roleID)
	if len(appRoles) == 0 {
		delete(s.roles, appID)
	}

	guildRoles, ok := s.rolesByGuild[appID]
	if ok {
		roles, ok := guildRoles[guildID]
		if ok {
			delete(roles, roleID)
			if len(roles) == 0 {
				delete(guildRoles, guildID)
				if len(guildRoles) == 0 {
					delete(s.rolesByGuild, appID)
				}
			}
		}
	}

	return nil
}

// CacheChannelStore methods

func (s *MapCacheStore) GetGuildChannel(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, channelID snowflake.ID) (*model.Channel, error) {
	s.channelsMu.RLock()
	defer s.channelsMu.RUnlock()

	guildChannels, ok := s.channelsByGuild[appID]
	if !ok {
		return nil, store.ErrNotFound
	}

	channels, ok := guildChannels[guildID]
	if !ok {
		return nil, store.ErrNotFound
	}

	channel, ok := channels[channelID]
	if !ok {
		return nil, store.ErrNotFound
	}

	return channel, nil
}

func (s *MapCacheStore) GetChannel(ctx context.Context, appID snowflake.ID, channelID snowflake.ID) (*model.Channel, error) {
	s.channelsMu.RLock()
	defer s.channelsMu.RUnlock()

	appChannels, ok := s.channels[appID]
	if !ok {
		return nil, store.ErrNotFound
	}

	channel, ok := appChannels[channelID]
	if !ok {
		return nil, store.ErrNotFound
	}

	return channel, nil
}

func (s *MapCacheStore) GetGuildChannels(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, limit int, offset int) ([]*model.Channel, error) {
	s.channelsMu.RLock()
	defer s.channelsMu.RUnlock()

	guildChannels, ok := s.channelsByGuild[appID]
	if !ok {
		return []*model.Channel{}, nil
	}

	channels, ok := guildChannels[guildID]
	if !ok {
		return []*model.Channel{}, nil
	}

	result := make([]*model.Channel, 0)
	if limit > 0 {
		result = make([]*model.Channel, 0, limit)
	}

	currentOffset := 0
	for _, channel := range channels {
		if currentOffset < offset {
			currentOffset++
			continue
		}

		result = append(result, channel)

		if limit > 0 && len(result) >= limit {
			break
		}
	}

	return result, nil
}

func (s *MapCacheStore) GetChannels(ctx context.Context, appID snowflake.ID, limit int, offset int) ([]*model.Channel, error) {
	s.channelsMu.RLock()
	defer s.channelsMu.RUnlock()

	appChannels, ok := s.channels[appID]
	if !ok {
		return []*model.Channel{}, nil
	}

	channels := make([]*model.Channel, 0)
	if limit > 0 {
		channels = make([]*model.Channel, 0, limit)
	}

	currentOffset := 0
	for _, channel := range appChannels {
		if currentOffset < offset {
			currentOffset++
			continue
		}

		channels = append(channels, channel)

		if limit > 0 && len(channels) >= limit {
			break
		}
	}

	return channels, nil
}

func (s *MapCacheStore) GetChannelsByType(ctx context.Context, appID snowflake.ID, types []int, limit int, offset int) ([]*model.Channel, error) {
	s.channelsMu.RLock()
	defer s.channelsMu.RUnlock()

	appChannels, ok := s.channels[appID]
	if !ok {
		return []*model.Channel{}, nil
	}

	channels := make([]*model.Channel, 0)
	if limit > 0 {
		channels = make([]*model.Channel, 0, limit)
	}

	currentOffset := 0
	for _, channel := range appChannels {
		if !slices.Contains(types, int(channel.Data.Type())) {
			continue
		}

		if currentOffset < offset {
			currentOffset++
			continue
		}

		channels = append(channels, channel)

		if limit > 0 && len(channels) >= limit {
			break
		}
	}

	return channels, nil
}

func (s *MapCacheStore) SearchGuildChannels(ctx context.Context, params store.SearchGuildChannelsParams) ([]*model.Channel, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *MapCacheStore) SearchChannels(ctx context.Context, params store.SearchChannelsParams) ([]*model.Channel, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *MapCacheStore) CountGuildChannels(ctx context.Context, appID snowflake.ID, guildID snowflake.ID) (int, error) {
	s.channelsMu.RLock()
	defer s.channelsMu.RUnlock()

	guildChannels, ok := s.channelsByGuild[appID]
	if !ok {
		return 0, nil
	}

	channels, ok := guildChannels[guildID]
	if !ok {
		return 0, nil
	}

	return len(channels), nil
}

func (s *MapCacheStore) CountChannels(ctx context.Context, appID snowflake.ID) (int, error) {
	s.channelsMu.RLock()
	defer s.channelsMu.RUnlock()

	appChannels, ok := s.channels[appID]
	if !ok {
		return 0, nil
	}

	return len(appChannels), nil
}

func (s *MapCacheStore) UpsertChannels(ctx context.Context, channels ...store.UpsertChannelParams) error {
	s.channelsMu.Lock()
	defer s.channelsMu.Unlock()

	for _, channel := range channels {
		if s.channels[channel.AppID] == nil {
			s.channels[channel.AppID] = make(map[snowflake.ID]*model.Channel)
		}
		if s.channelsByGuild[channel.AppID] == nil {
			s.channelsByGuild[channel.AppID] = make(map[snowflake.ID]map[snowflake.ID]*model.Channel)
		}
		if s.channelsByGuild[channel.AppID][channel.GuildID] == nil {
			s.channelsByGuild[channel.AppID][channel.GuildID] = make(map[snowflake.ID]*model.Channel)
		}

		createdAt := channel.CreatedAt
		if createdAt.IsZero() {
			createdAt = time.Now().UTC()
		}
		updatedAt := channel.UpdatedAt
		if updatedAt.IsZero() {
			updatedAt = time.Now().UTC()
		}

		c := &model.Channel{
			AppID:     channel.AppID,
			GuildID:   channel.GuildID,
			ChannelID: channel.ChannelID,
			Data:      channel.Data,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}

		s.channels[channel.AppID][channel.ChannelID] = c
		s.channelsByGuild[channel.AppID][channel.GuildID][channel.ChannelID] = c
	}

	return nil
}

func (s *MapCacheStore) DeleteChannel(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, channelID snowflake.ID) error {
	s.channelsMu.Lock()
	defer s.channelsMu.Unlock()

	appChannels, ok := s.channels[appID]
	if !ok {
		return nil
	}

	_, ok = appChannels[channelID]
	if !ok {
		return nil
	}

	delete(appChannels, channelID)
	if len(appChannels) == 0 {
		delete(s.channels, appID)
	}

	guildChannels, ok := s.channelsByGuild[appID]
	if ok {
		channels, ok := guildChannels[guildID]
		if ok {
			delete(channels, channelID)
			if len(channels) == 0 {
				delete(guildChannels, guildID)
				if len(guildChannels) == 0 {
					delete(s.channelsByGuild, appID)
				}
			}
		}
	}

	return nil
}

// CacheEmojiStore methods

func (s *MapCacheStore) GetGuildEmoji(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, emojiID snowflake.ID) (*model.Emoji, error) {
	s.emojisMu.RLock()
	defer s.emojisMu.RUnlock()

	guildEmojis, ok := s.emojisByGuild[appID]
	if !ok {
		return nil, store.ErrNotFound
	}

	emojis, ok := guildEmojis[guildID]
	if !ok {
		return nil, store.ErrNotFound
	}

	emoji, ok := emojis[emojiID]
	if !ok {
		return nil, store.ErrNotFound
	}

	return emoji, nil
}

func (s *MapCacheStore) GetEmoji(ctx context.Context, appID snowflake.ID, emojiID snowflake.ID) (*model.Emoji, error) {
	s.emojisMu.RLock()
	defer s.emojisMu.RUnlock()

	appEmojis, ok := s.emojis[appID]
	if !ok {
		return nil, store.ErrNotFound
	}

	emoji, ok := appEmojis[emojiID]
	if !ok {
		return nil, store.ErrNotFound
	}

	return emoji, nil
}

func (s *MapCacheStore) GetGuildEmojis(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, limit int, offset int) ([]*model.Emoji, error) {
	s.emojisMu.RLock()
	defer s.emojisMu.RUnlock()

	guildEmojis, ok := s.emojisByGuild[appID]
	if !ok {
		return []*model.Emoji{}, nil
	}

	emojis, ok := guildEmojis[guildID]
	if !ok {
		return []*model.Emoji{}, nil
	}

	result := make([]*model.Emoji, 0)
	if limit > 0 {
		result = make([]*model.Emoji, 0, limit)
	}

	currentOffset := 0
	for _, emoji := range emojis {
		if currentOffset < offset {
			currentOffset++
			continue
		}

		result = append(result, emoji)

		if limit > 0 && len(result) >= limit {
			break
		}
	}

	return result, nil
}

func (s *MapCacheStore) GetEmojis(ctx context.Context, appID snowflake.ID, limit int, offset int) ([]*model.Emoji, error) {
	s.emojisMu.RLock()
	defer s.emojisMu.RUnlock()

	appEmojis, ok := s.emojis[appID]
	if !ok {
		return []*model.Emoji{}, nil
	}

	emojis := make([]*model.Emoji, 0)
	if limit > 0 {
		emojis = make([]*model.Emoji, 0, limit)
	}

	currentOffset := 0
	for _, emoji := range appEmojis {
		if currentOffset < offset {
			currentOffset++
			continue
		}

		emojis = append(emojis, emoji)

		if limit > 0 && len(emojis) >= limit {
			break
		}
	}

	return emojis, nil
}

func (s *MapCacheStore) SearchGuildEmojis(ctx context.Context, params store.SearchGuildEmojisParams) ([]*model.Emoji, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *MapCacheStore) SearchEmojis(ctx context.Context, params store.SearchEmojisParams) ([]*model.Emoji, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *MapCacheStore) CountGuildEmojis(ctx context.Context, appID snowflake.ID, guildID snowflake.ID) (int, error) {
	s.emojisMu.RLock()
	defer s.emojisMu.RUnlock()

	guildEmojis, ok := s.emojisByGuild[appID]
	if !ok {
		return 0, nil
	}

	emojis, ok := guildEmojis[guildID]
	if !ok {
		return 0, nil
	}

	return len(emojis), nil
}

func (s *MapCacheStore) CountEmojis(ctx context.Context, appID snowflake.ID) (int, error) {
	s.emojisMu.RLock()
	defer s.emojisMu.RUnlock()

	appEmojis, ok := s.emojis[appID]
	if !ok {
		return 0, nil
	}

	return len(appEmojis), nil
}

func (s *MapCacheStore) UpsertEmojis(ctx context.Context, emojis ...store.UpsertEmojiParams) error {
	s.emojisMu.Lock()
	defer s.emojisMu.Unlock()

	for _, emoji := range emojis {
		if s.emojis[emoji.AppID] == nil {
			s.emojis[emoji.AppID] = make(map[snowflake.ID]*model.Emoji)
		}
		if s.emojisByGuild[emoji.AppID] == nil {
			s.emojisByGuild[emoji.AppID] = make(map[snowflake.ID]map[snowflake.ID]*model.Emoji)
		}
		if s.emojisByGuild[emoji.AppID][emoji.GuildID] == nil {
			s.emojisByGuild[emoji.AppID][emoji.GuildID] = make(map[snowflake.ID]*model.Emoji)
		}

		createdAt := emoji.CreatedAt
		if createdAt.IsZero() {
			createdAt = time.Now().UTC()
		}
		updatedAt := emoji.UpdatedAt
		if updatedAt.IsZero() {
			updatedAt = time.Now().UTC()
		}

		e := &model.Emoji{
			AppID:     emoji.AppID,
			GuildID:   emoji.GuildID,
			EmojiID:   emoji.EmojiID,
			Data:      emoji.Data,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}

		s.emojis[emoji.AppID][emoji.EmojiID] = e
		s.emojisByGuild[emoji.AppID][emoji.GuildID][emoji.EmojiID] = e
	}

	return nil
}

func (s *MapCacheStore) DeleteEmoji(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, emojiID snowflake.ID) error {
	s.emojisMu.Lock()
	defer s.emojisMu.Unlock()

	appEmojis, ok := s.emojis[appID]
	if !ok {
		return nil
	}

	_, ok = appEmojis[emojiID]
	if !ok {
		return nil
	}

	delete(appEmojis, emojiID)
	if len(appEmojis) == 0 {
		delete(s.emojis, appID)
	}

	guildEmojis, ok := s.emojisByGuild[appID]
	if ok {
		emojis, ok := guildEmojis[guildID]
		if ok {
			delete(emojis, emojiID)
			if len(emojis) == 0 {
				delete(guildEmojis, guildID)
				if len(guildEmojis) == 0 {
					delete(s.emojisByGuild, appID)
				}
			}
		}
	}

	return nil
}

// CacheStickerStore methods

func (s *MapCacheStore) GetSticker(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, stickerID snowflake.ID) (*model.Sticker, error) {
	s.stickersMu.RLock()
	defer s.stickersMu.RUnlock()

	guildStickers, ok := s.stickersByGuild[appID]
	if !ok {
		return nil, store.ErrNotFound
	}

	stickers, ok := guildStickers[guildID]
	if !ok {
		return nil, store.ErrNotFound
	}

	sticker, ok := stickers[stickerID]
	if !ok {
		return nil, store.ErrNotFound
	}

	return sticker, nil
}

func (s *MapCacheStore) GetGuildSticker(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, stickerID snowflake.ID) (*model.Sticker, error) {
	return s.GetSticker(ctx, appID, guildID, stickerID)
}

func (s *MapCacheStore) GetStickers(ctx context.Context, appID snowflake.ID, limit int, offset int) ([]*model.Sticker, error) {
	s.stickersMu.RLock()
	defer s.stickersMu.RUnlock()

	appStickers, ok := s.stickers[appID]
	if !ok {
		return []*model.Sticker{}, nil
	}

	stickers := make([]*model.Sticker, 0)
	if limit > 0 {
		stickers = make([]*model.Sticker, 0, limit)
	}

	currentOffset := 0
	for _, sticker := range appStickers {
		if currentOffset < offset {
			currentOffset++
			continue
		}

		stickers = append(stickers, sticker)

		if limit > 0 && len(stickers) >= limit {
			break
		}
	}

	return stickers, nil
}

func (s *MapCacheStore) GetGuildStickers(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, limit int, offset int) ([]*model.Sticker, error) {
	s.stickersMu.RLock()
	defer s.stickersMu.RUnlock()

	guildStickers, ok := s.stickersByGuild[appID]
	if !ok {
		return []*model.Sticker{}, nil
	}

	stickers, ok := guildStickers[guildID]
	if !ok {
		return []*model.Sticker{}, nil
	}

	result := make([]*model.Sticker, 0)
	if limit > 0 {
		result = make([]*model.Sticker, 0, limit)
	}

	currentOffset := 0
	for _, sticker := range stickers {
		if currentOffset < offset {
			currentOffset++
			continue
		}

		result = append(result, sticker)

		if limit > 0 && len(result) >= limit {
			break
		}
	}

	return result, nil
}

func (s *MapCacheStore) SearchStickers(ctx context.Context, params store.SearchStickersParams) ([]*model.Sticker, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *MapCacheStore) SearchGuildStickers(ctx context.Context, params store.SearchGuildStickersParams) ([]*model.Sticker, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *MapCacheStore) CountGuildStickers(ctx context.Context, appID snowflake.ID, guildID snowflake.ID) (int, error) {
	s.stickersMu.RLock()
	defer s.stickersMu.RUnlock()

	guildStickers, ok := s.stickersByGuild[appID]
	if !ok {
		return 0, nil
	}

	stickers, ok := guildStickers[guildID]
	if !ok {
		return 0, nil
	}

	return len(stickers), nil
}

func (s *MapCacheStore) CountStickers(ctx context.Context, appID snowflake.ID) (int, error) {
	s.stickersMu.RLock()
	defer s.stickersMu.RUnlock()

	appStickers, ok := s.stickers[appID]
	if !ok {
		return 0, nil
	}

	return len(appStickers), nil
}

func (s *MapCacheStore) UpsertStickers(ctx context.Context, stickers ...store.UpsertStickerParams) error {
	s.stickersMu.Lock()
	defer s.stickersMu.Unlock()

	for _, sticker := range stickers {
		if s.stickers[sticker.AppID] == nil {
			s.stickers[sticker.AppID] = make(map[snowflake.ID]*model.Sticker)
		}
		if s.stickersByGuild[sticker.AppID] == nil {
			s.stickersByGuild[sticker.AppID] = make(map[snowflake.ID]map[snowflake.ID]*model.Sticker)
		}
		if s.stickersByGuild[sticker.AppID][sticker.GuildID] == nil {
			s.stickersByGuild[sticker.AppID][sticker.GuildID] = make(map[snowflake.ID]*model.Sticker)
		}

		createdAt := sticker.CreatedAt
		if createdAt.IsZero() {
			createdAt = time.Now().UTC()
		}
		updatedAt := sticker.UpdatedAt
		if updatedAt.IsZero() {
			updatedAt = time.Now().UTC()
		}

		st := &model.Sticker{
			AppID:     sticker.AppID,
			GuildID:   sticker.GuildID,
			StickerID: sticker.StickerID,
			Data:      sticker.Data,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}

		s.stickers[sticker.AppID][sticker.StickerID] = st
		s.stickersByGuild[sticker.AppID][sticker.GuildID][sticker.StickerID] = st
	}

	return nil
}

func (s *MapCacheStore) DeleteSticker(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, stickerID snowflake.ID) error {
	s.stickersMu.Lock()
	defer s.stickersMu.Unlock()

	appStickers, ok := s.stickers[appID]
	if !ok {
		return nil
	}

	_, ok = appStickers[stickerID]
	if !ok {
		return nil
	}

	delete(appStickers, stickerID)
	if len(appStickers) == 0 {
		delete(s.stickers, appID)
	}

	guildStickers, ok := s.stickersByGuild[appID]
	if ok {
		stickers, ok := guildStickers[guildID]
		if ok {
			delete(stickers, stickerID)
			if len(stickers) == 0 {
				delete(guildStickers, guildID)
				if len(guildStickers) == 0 {
					delete(s.stickersByGuild, appID)
				}
			}
		}
	}

	return nil
}

// CacheStore methods

func (s *MapCacheStore) MarkShardEntitiesTainted(ctx context.Context, params store.MarkShardEntitiesTaintedParams) error {
	// Not implemented for in-memory store
	return nil
}

func (s *MapCacheStore) MassUpsertEntities(ctx context.Context, params store.MassUpsertEntitiesParams) error {
	if len(params.Guilds) > 0 {
		if err := s.UpsertGuilds(ctx, params.Guilds...); err != nil {
			return err
		}
	}

	if len(params.Roles) > 0 {
		if err := s.UpsertRoles(ctx, params.Roles...); err != nil {
			return err
		}
	}

	if len(params.Channels) > 0 {
		if err := s.UpsertChannels(ctx, params.Channels...); err != nil {
			return err
		}
	}

	if len(params.Emojis) > 0 {
		if err := s.UpsertEmojis(ctx, params.Emojis...); err != nil {
			return err
		}
	}

	if len(params.Stickers) > 0 {
		if err := s.UpsertStickers(ctx, params.Stickers...); err != nil {
			return err
		}
	}

	return nil
}
