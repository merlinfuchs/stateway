package server

import (
	"context"
	"slices"
	"sync"
	"time"

	discache "github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/snowflake/v2"
	"github.com/merlinfuchs/stateway/stateway-cache/model"
	"github.com/merlinfuchs/stateway/stateway-cache/store"
	"github.com/merlinfuchs/stateway/stateway-lib/cache"
)

var _ store.CacheStore = (*InMemoryCacheStore)(nil)

type AppCaches struct {
	Guilds   discache.Cache[*cache.Guild]
	Channels discache.Cache[*cache.Channel]
	Roles    discache.Cache[*cache.Role]
	Emojis   discache.Cache[*cache.Emoji]
	Stickers discache.Cache[*cache.Sticker]
}

func NewAppCaches() *AppCaches {
	return &AppCaches{
		Guilds:   discache.NewCache[*cache.Guild](0, 0, discache.PolicyAll),
		Channels: discache.NewCache[*cache.Channel](0, 0, discache.PolicyAll),
		Roles:    discache.NewCache[*cache.Role](0, 0, discache.PolicyAll),
		Emojis:   discache.NewCache[*cache.Emoji](0, 0, discache.PolicyAll),
		Stickers: discache.NewCache[*cache.Sticker](0, 0, discache.PolicyAll),
	}
}

type InMemoryCacheStore struct {
	sync.RWMutex
	createOpts []discache.ConfigOpt
	caches     map[snowflake.ID]*AppCaches
}

func NewInMemoryCacheStore(opts ...discache.ConfigOpt) *InMemoryCacheStore {
	return &InMemoryCacheStore{
		createOpts: opts,
		caches:     make(map[snowflake.ID]*AppCaches),
	}
}

func (s *InMemoryCacheStore) GetCaches(appID snowflake.ID) *AppCaches {
	s.RLock()
	caches, ok := s.caches[appID]
	s.RUnlock()

	if !ok {
		caches = NewAppCaches()
		s.Lock()
		s.caches[appID] = caches
		s.Unlock()
	}

	return caches
}

// CacheGuildStore methods

func (s *InMemoryCacheStore) GetGuild(ctx context.Context, appID snowflake.ID, guildID snowflake.ID) (*model.Guild, error) {
	caches := s.GetCaches(appID)
	guild, ok := caches.Guilds.Get(guildID)
	if !ok {
		return nil, store.ErrNotFound
	}
	return guild, nil
}

func (s *InMemoryCacheStore) GetGuildOwnerID(ctx context.Context, appID snowflake.ID, guildID snowflake.ID) (snowflake.ID, error) {
	caches := s.GetCaches(appID)
	guild, ok := caches.Guilds.Get(guildID)
	if !ok {
		return 0, store.ErrNotFound
	}
	return guild.Data.OwnerID, nil
}

func (s *InMemoryCacheStore) GetGuilds(ctx context.Context, appID snowflake.ID, limit int, offset int) ([]*model.Guild, error) {
	caches := s.GetCaches(appID)
	guilds := make([]*model.Guild, 0)
	for guild := range caches.Guilds.All() {
		guilds = append(guilds, guild)
	}
	return guilds, nil
}

func (s *InMemoryCacheStore) CheckGuildExist(ctx context.Context, appID snowflake.ID, guildID snowflake.ID) (bool, error) {
	_, ok := s.GetCaches(appID).Guilds.Get(guildID)
	return ok, nil
}

func (s *InMemoryCacheStore) UpsertGuilds(ctx context.Context, guilds ...store.UpsertGuildParams) error {
	for _, guild := range guilds {
		caches := s.GetCaches(guild.AppID)
		caches.Guilds.Put(guild.GuildID, &model.Guild{
			AppID:     guild.AppID,
			GuildID:   guild.GuildID,
			Data:      guild.Data,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
	}
	return nil
}

func (s *InMemoryCacheStore) MarkGuildUnavailable(ctx context.Context, appID snowflake.ID, guildID snowflake.ID) error {
	caches := s.GetCaches(appID)
	guild, ok := caches.Guilds.Get(guildID)
	if !ok {
		return store.ErrNotFound
	}
	guild.Unavailable = true
	guild.UpdatedAt = time.Now().UTC()
	caches.Guilds.Put(guildID, guild)
	return nil
}

func (s *InMemoryCacheStore) DeleteGuild(ctx context.Context, appID snowflake.ID, guildID snowflake.ID) error {
	caches := s.GetCaches(appID)
	caches.Guilds.Remove(guildID)
	return nil
}

func (s *InMemoryCacheStore) SearchGuilds(ctx context.Context, params store.SearchGuildsParams) ([]*model.Guild, error) {
	return nil, nil
}

// CacheRoleStore methods

func (s *InMemoryCacheStore) GetGuildRole(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, roleID snowflake.ID) (*model.Role, error) {
	caches := s.GetCaches(appID)
	role, ok := caches.Roles.Get(roleID)
	if !ok {
		return nil, store.ErrNotFound
	}
	if role.GuildID != guildID {
		return nil, store.ErrNotFound
	}
	return role, nil
}

func (s *InMemoryCacheStore) GetRole(ctx context.Context, appID snowflake.ID, roleID snowflake.ID) (*model.Role, error) {
	caches := s.GetCaches(appID)
	role, ok := caches.Roles.Get(roleID)
	if !ok {
		return nil, store.ErrNotFound
	}
	return role, nil
}

func (s *InMemoryCacheStore) GetGuildRoles(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, limit int, offset int) ([]*model.Role, error) {
	caches := s.GetCaches(appID)
	roles := make([]*model.Role, 0)
	for role := range caches.Roles.All() {
		if role.GuildID != guildID {
			continue
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (s *InMemoryCacheStore) GetGuildRolesByIDs(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, roleIDs []snowflake.ID) ([]*model.Role, error) {
	caches := s.GetCaches(appID)
	roles := make([]*model.Role, 0)
	for _, roleID := range roleIDs {
		role, ok := caches.Roles.Get(roleID)
		if !ok || role.GuildID != guildID {
			continue
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (s *InMemoryCacheStore) GetRoles(ctx context.Context, appID snowflake.ID, limit int, offset int) ([]*model.Role, error) {
	caches := s.GetCaches(appID)
	roles := make([]*model.Role, 0)
	for role := range caches.Roles.All() {
		roles = append(roles, role)
	}
	return roles, nil
}

func (s *InMemoryCacheStore) SearchGuildRoles(ctx context.Context, params store.SearchGuildRolesParams) ([]*model.Role, error) {
	return nil, nil
}

func (s *InMemoryCacheStore) SearchRoles(ctx context.Context, params store.SearchRolesParams) ([]*model.Role, error) {
	return nil, nil
}

func (s *InMemoryCacheStore) CountGuildRoles(ctx context.Context, appID snowflake.ID, guildID snowflake.ID) (int, error) {
	caches := s.GetCaches(appID)
	var count int
	for role := range caches.Roles.All() {
		if role.GuildID == guildID {
			count++
		}
	}
	return count, nil
}

func (s *InMemoryCacheStore) CountRoles(ctx context.Context, appID snowflake.ID) (int, error) {
	caches := s.GetCaches(appID)
	return caches.Roles.Len(), nil
}

func (s *InMemoryCacheStore) UpsertRoles(ctx context.Context, roles ...store.UpsertRoleParams) error {
	for _, role := range roles {
		caches := s.GetCaches(role.AppID)
		caches.Roles.Put(role.RoleID, &model.Role{
			AppID:   role.AppID,
			GuildID: role.GuildID,
			RoleID:  role.RoleID,
			Data:    role.Data,
		})
	}
	return nil
}

func (s *InMemoryCacheStore) DeleteRole(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, roleID snowflake.ID) error {
	caches := s.GetCaches(appID)
	caches.Roles.Remove(roleID)
	return nil
}

// CacheChannelStore methods

func (s *InMemoryCacheStore) GetGuildChannel(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, channelID snowflake.ID) (*model.Channel, error) {
	caches := s.GetCaches(appID)
	channel, ok := caches.Channels.Get(channelID)
	if !ok {
		return nil, store.ErrNotFound
	}
	if channel.GuildID != guildID {
		return nil, store.ErrNotFound
	}
	return channel, nil
}

func (s *InMemoryCacheStore) GetChannel(ctx context.Context, appID snowflake.ID, channelID snowflake.ID) (*model.Channel, error) {
	caches := s.GetCaches(appID)
	channel, ok := caches.Channels.Get(channelID)
	if !ok {
		return nil, store.ErrNotFound
	}
	return channel, nil
}

func (s *InMemoryCacheStore) GetGuildChannels(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, limit int, offset int) ([]*model.Channel, error) {
	caches := s.GetCaches(appID)
	channels := make([]*model.Channel, 0)
	for channel := range caches.Channels.All() {
		if channel.GuildID != guildID {
			continue
		}
		channels = append(channels, channel)
	}
	return channels, nil
}

func (s *InMemoryCacheStore) GetChannels(ctx context.Context, appID snowflake.ID, limit int, offset int) ([]*model.Channel, error) {
	caches := s.GetCaches(appID)
	channels := make([]*model.Channel, 0)
	for channel := range caches.Channels.All() {
		channels = append(channels, channel)
	}
	return channels, nil
}

func (s *InMemoryCacheStore) GetChannelsByType(ctx context.Context, appID snowflake.ID, types []int, limit int, offset int) ([]*model.Channel, error) {
	caches := s.GetCaches(appID)
	channels := make([]*model.Channel, 0)
	for channel := range caches.Channels.All() {
		if !slices.Contains(types, int(channel.Data.Type())) {
			continue
		}
		channels = append(channels, channel)
	}
	return channels, nil
}

func (s *InMemoryCacheStore) SearchGuildChannels(ctx context.Context, params store.SearchGuildChannelsParams) ([]*model.Channel, error) {
	return nil, nil
}

func (s *InMemoryCacheStore) SearchChannels(ctx context.Context, params store.SearchChannelsParams) ([]*model.Channel, error) {
	return nil, nil
}

func (s *InMemoryCacheStore) CountGuildChannels(ctx context.Context, appID snowflake.ID, guildID snowflake.ID) (int, error) {
	caches := s.GetCaches(appID)
	var count int
	for channel := range caches.Channels.All() {
		if channel.GuildID == guildID {
			count++
		}
	}
	return count, nil
}

func (s *InMemoryCacheStore) CountChannels(ctx context.Context, appID snowflake.ID) (int, error) {
	caches := s.GetCaches(appID)
	return caches.Channels.Len(), nil
}

func (s *InMemoryCacheStore) UpsertChannels(ctx context.Context, channels ...store.UpsertChannelParams) error {
	for _, channel := range channels {
		caches := s.GetCaches(channel.AppID)
		caches.Channels.Put(channel.ChannelID, &model.Channel{
			AppID:     channel.AppID,
			GuildID:   channel.GuildID,
			ChannelID: channel.ChannelID,
			Data:      channel.Data,
		})
	}
	return nil
}

func (s *InMemoryCacheStore) DeleteChannel(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, channelID snowflake.ID) error {
	caches := s.GetCaches(appID)
	caches.Channels.Remove(channelID)
	return nil
}

// CacheEmojiStore methods

func (s *InMemoryCacheStore) GetGuildEmoji(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, emojiID snowflake.ID) (*model.Emoji, error) {
	caches := s.GetCaches(appID)
	emoji, ok := caches.Emojis.Get(emojiID)
	if !ok {
		return nil, store.ErrNotFound
	}
	if emoji.GuildID != guildID {
		return nil, store.ErrNotFound
	}
	return emoji, nil
}

func (s *InMemoryCacheStore) GetEmoji(ctx context.Context, appID snowflake.ID, emojiID snowflake.ID) (*model.Emoji, error) {
	caches := s.GetCaches(appID)
	emoji, ok := caches.Emojis.Get(emojiID)
	if !ok {
		return nil, store.ErrNotFound
	}
	return emoji, nil
}

func (s *InMemoryCacheStore) GetGuildEmojis(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, limit int, offset int) ([]*model.Emoji, error) {
	caches := s.GetCaches(appID)
	emojis := make([]*model.Emoji, 0)
	for emoji := range caches.Emojis.All() {
		if emoji.GuildID != guildID {
			continue
		}
		emojis = append(emojis, emoji)
	}
	return emojis, nil
}

func (s *InMemoryCacheStore) GetEmojis(ctx context.Context, appID snowflake.ID, limit int, offset int) ([]*model.Emoji, error) {
	caches := s.GetCaches(appID)
	emojis := make([]*model.Emoji, 0)
	for emoji := range caches.Emojis.All() {
		emojis = append(emojis, emoji)
	}
	return emojis, nil
}

func (s *InMemoryCacheStore) SearchGuildEmojis(ctx context.Context, params store.SearchGuildEmojisParams) ([]*model.Emoji, error) {
	return nil, nil
}

func (s *InMemoryCacheStore) SearchEmojis(ctx context.Context, params store.SearchEmojisParams) ([]*model.Emoji, error) {
	return nil, nil
}

func (s *InMemoryCacheStore) CountGuildEmojis(ctx context.Context, appID snowflake.ID, guildID snowflake.ID) (int, error) {
	caches := s.GetCaches(appID)
	var count int
	for emoji := range caches.Emojis.All() {
		if emoji.GuildID == guildID {
			count++
		}
	}
	return count, nil
}

func (s *InMemoryCacheStore) CountEmojis(ctx context.Context, appID snowflake.ID) (int, error) {
	caches := s.GetCaches(appID)
	return caches.Emojis.Len(), nil
}

func (s *InMemoryCacheStore) UpsertEmojis(ctx context.Context, emojis ...store.UpsertEmojiParams) error {
	for _, emoji := range emojis {
		caches := s.GetCaches(emoji.AppID)
		caches.Emojis.Put(emoji.EmojiID, &model.Emoji{
			AppID:   emoji.AppID,
			GuildID: emoji.GuildID,
			EmojiID: emoji.EmojiID,
			Data:    emoji.Data,
		})
	}
	return nil
}

func (s *InMemoryCacheStore) DeleteEmoji(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, emojiID snowflake.ID) error {
	caches := s.GetCaches(appID)
	caches.Emojis.Remove(emojiID)
	return nil
}

// CacheStickerStore methods

func (s *InMemoryCacheStore) GetSticker(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, stickerID snowflake.ID) (*model.Sticker, error) {
	caches := s.GetCaches(appID)
	sticker, ok := caches.Stickers.Get(stickerID)
	if !ok {
		return nil, store.ErrNotFound
	}
	if sticker.GuildID != guildID {
		return nil, store.ErrNotFound
	}
	return sticker, nil
}

func (s *InMemoryCacheStore) GetGuildSticker(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, stickerID snowflake.ID) (*model.Sticker, error) {
	caches := s.GetCaches(appID)
	sticker, ok := caches.Stickers.Get(stickerID)
	if !ok {
		return nil, store.ErrNotFound
	}
	if sticker.GuildID != guildID {
		return nil, store.ErrNotFound
	}
	return sticker, nil
}

func (s *InMemoryCacheStore) GetStickers(ctx context.Context, appID snowflake.ID, limit int, offset int) ([]*model.Sticker, error) {
	caches := s.GetCaches(appID)
	stickers := make([]*model.Sticker, 0)
	for sticker := range caches.Stickers.All() {
		stickers = append(stickers, sticker)
	}
	return stickers, nil
}

func (s *InMemoryCacheStore) GetGuildStickers(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, limit int, offset int) ([]*model.Sticker, error) {
	caches := s.GetCaches(appID)
	stickers := make([]*model.Sticker, 0)
	for sticker := range caches.Stickers.All() {
		if sticker.GuildID != guildID {
			continue
		}
		stickers = append(stickers, sticker)
	}
	return stickers, nil
}

func (s *InMemoryCacheStore) SearchStickers(ctx context.Context, params store.SearchStickersParams) ([]*model.Sticker, error) {
	return nil, nil
}

func (s *InMemoryCacheStore) SearchGuildStickers(ctx context.Context, params store.SearchGuildStickersParams) ([]*model.Sticker, error) {
	return nil, nil
}

func (s *InMemoryCacheStore) CountGuildStickers(ctx context.Context, appID snowflake.ID, guildID snowflake.ID) (int, error) {
	caches := s.GetCaches(appID)
	var count int
	for sticker := range caches.Stickers.All() {
		if sticker.GuildID == guildID {
			count++
		}
	}
	return count, nil
}

func (s *InMemoryCacheStore) CountStickers(ctx context.Context, appID snowflake.ID) (int, error) {
	caches := s.GetCaches(appID)
	return caches.Stickers.Len(), nil
}

func (s *InMemoryCacheStore) UpsertStickers(ctx context.Context, stickers ...store.UpsertStickerParams) error {
	for _, sticker := range stickers {
		caches := s.GetCaches(sticker.AppID)
		caches.Stickers.Put(sticker.StickerID, &model.Sticker{
			AppID:     sticker.AppID,
			GuildID:   sticker.GuildID,
			StickerID: sticker.StickerID,
			Data:      sticker.Data,
		})
	}
	return nil
}

func (s *InMemoryCacheStore) DeleteSticker(ctx context.Context, appID snowflake.ID, guildID snowflake.ID, stickerID snowflake.ID) error {
	caches := s.GetCaches(appID)
	caches.Stickers.Remove(stickerID)
	return nil
}

// CacheStore methods

func (s *InMemoryCacheStore) MarkShardEntitiesTainted(ctx context.Context, params store.MarkShardEntitiesTaintedParams) error {
	return nil
}

func (s *InMemoryCacheStore) MassUpsertEntities(ctx context.Context, params store.MassUpsertEntitiesParams) error {
	caches := s.GetCaches(params.AppID)
	for _, guild := range params.Guilds {
		caches.Guilds.Put(guild.GuildID, &model.Guild{
			AppID:     guild.AppID,
			GuildID:   guild.GuildID,
			Data:      guild.Data,
			CreatedAt: guild.CreatedAt,
			UpdatedAt: guild.UpdatedAt,
		})
	}
	for _, role := range params.Roles {
		caches.Roles.Put(role.RoleID, &model.Role{
			AppID:     role.AppID,
			GuildID:   role.GuildID,
			RoleID:    role.RoleID,
			Data:      role.Data,
			CreatedAt: role.CreatedAt,
			UpdatedAt: role.UpdatedAt,
		})
	}
	for _, channel := range params.Channels {
		caches.Channels.Put(channel.ChannelID, &model.Channel{
			AppID:     channel.AppID,
			GuildID:   channel.GuildID,
			ChannelID: channel.ChannelID,
			Data:      channel.Data,
			CreatedAt: channel.CreatedAt,
			UpdatedAt: channel.UpdatedAt,
		})
	}
	for _, emoji := range params.Emojis {
		caches.Emojis.Put(emoji.EmojiID, &model.Emoji{
			AppID:     emoji.AppID,
			GuildID:   emoji.GuildID,
			EmojiID:   emoji.EmojiID,
			Data:      emoji.Data,
			CreatedAt: emoji.CreatedAt,
			UpdatedAt: emoji.UpdatedAt,
		})
	}
	for _, sticker := range params.Stickers {
		caches.Stickers.Put(sticker.StickerID, &model.Sticker{
			AppID:     sticker.AppID,
			GuildID:   sticker.GuildID,
			StickerID: sticker.StickerID,
			Data:      sticker.Data,
		})
	}
	return nil
}
