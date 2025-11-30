package server

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
	"github.com/merlinfuchs/stateway/stateway-cache/store"
)

type GuildCreateBatcher struct {
	cacheStore store.CacheStore
	batchSize  int
	timeout    time.Duration

	mu         sync.Mutex
	batch      []pendingGuildCreate
	flushTimer *time.Timer
	ctx        context.Context
	cancel     context.CancelFunc
}

type pendingGuildCreate struct {
	appID   snowflake.ID
	guildID snowflake.ID
	event   gateway.EventGuildCreate
}

func NewGuildCreateBatcher(cacheStore store.CacheStore, batchSize int, timeout time.Duration) *GuildCreateBatcher {
	ctx, cancel := context.WithCancel(context.Background())
	return &GuildCreateBatcher{
		cacheStore: cacheStore,
		batchSize:  batchSize,
		timeout:    timeout,
		ctx:        ctx,
		cancel:     cancel,
	}
}

func (b *GuildCreateBatcher) Add(ctx context.Context, appID snowflake.ID, event gateway.EventGuildCreate) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.batch = append(b.batch, pendingGuildCreate{
		appID:   appID,
		guildID: event.ID,
		event:   event,
	})

	// If we've reached the batch size, flush immediately
	if len(b.batch) >= b.batchSize {
		return b.flushLocked(ctx)
	}

	// Reset the flush timer
	if b.flushTimer != nil {
		b.flushTimer.Stop()
	}
	b.flushTimer = time.AfterFunc(b.timeout, func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		if len(b.batch) > 0 {
			// Use a context with timeout for the flush operation
			flushCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			if err := b.flushLocked(flushCtx); err != nil {
				slog.Error("Failed to flush batched guild creates", slog.String("error", err.Error()))
			}
		}
	})

	return nil
}

func (b *GuildCreateBatcher) flushLocked(ctx context.Context) error {
	if len(b.batch) == 0 {
		return nil
	}

	// Stop the timer if it's running
	if b.flushTimer != nil {
		b.flushTimer.Stop()
		b.flushTimer = nil
	}

	// Prepare the batch
	batch := b.batch
	b.batch = nil

	// Release the lock before processing
	b.mu.Unlock()

	// Process the batch
	err := b.processBatch(ctx, batch)

	// Re-acquire the lock
	b.mu.Lock()

	return err
}

func (b *GuildCreateBatcher) processBatch(ctx context.Context, batch []pendingGuildCreate) error {
	// Aggregate all entities from all guilds
	var allGuilds []store.UpsertGuildParams
	var allRoles []store.UpsertRoleParams
	var allChannels []store.UpsertChannelParams
	var allEmojis []store.UpsertEmojiParams
	var allStickers []store.UpsertStickerParams

	now := time.Now().UTC()

	for _, pending := range batch {
		// Guild
		allGuilds = append(allGuilds, store.UpsertGuildParams{
			AppID:     pending.appID,
			GuildID:   pending.guildID,
			Data:      pending.event.Guild,
			CreatedAt: now,
			UpdatedAt: now,
		})

		// Roles
		for _, role := range pending.event.Roles {
			allRoles = append(allRoles, store.UpsertRoleParams{
				AppID:     pending.appID,
				GuildID:   pending.guildID,
				RoleID:    role.ID,
				Data:      role,
				CreatedAt: now,
				UpdatedAt: now,
			})
		}

		// Channels
		for _, channel := range pending.event.Channels {
			channel := ensureChannelGuildID(channel, pending.guildID)
			allChannels = append(allChannels, store.UpsertChannelParams{
				AppID:     pending.appID,
				GuildID:   pending.guildID,
				ChannelID: channel.ID(),
				Data:      channel,
				CreatedAt: now,
				UpdatedAt: now,
			})
		}
		for _, thread := range pending.event.Threads {
			channel := ensureChannelGuildID(thread, pending.guildID)
			allChannels = append(allChannels, store.UpsertChannelParams{
				AppID:     pending.appID,
				GuildID:   pending.guildID,
				ChannelID: thread.ID(),
				Data:      channel,
				CreatedAt: now,
				UpdatedAt: now,
			})
		}

		// Emojis
		for _, emoji := range pending.event.Emojis {
			allEmojis = append(allEmojis, store.UpsertEmojiParams{
				AppID:     pending.appID,
				GuildID:   pending.guildID,
				EmojiID:   emoji.ID,
				Data:      emoji,
				CreatedAt: now,
				UpdatedAt: now,
			})
		}

		// Stickers
		for _, sticker := range pending.event.Stickers {
			allStickers = append(allStickers, store.UpsertStickerParams{
				AppID:     pending.appID,
				GuildID:   pending.guildID,
				StickerID: sticker.ID,
				Data:      sticker,
				CreatedAt: now,
				UpdatedAt: now,
			})
		}
	}

	err := b.cacheStore.MassUpsertEntities(ctx, store.MassUpsertEntitiesParams{
		Guilds:   allGuilds,
		Roles:    allRoles,
		Channels: allChannels,
		Emojis:   allEmojis,
		Stickers: allStickers,
	})
	if err != nil {
		return fmt.Errorf("failed to mass upsert entities: %w", err)
	}

	slog.Debug(
		"Flushed batched guild creates",
		slog.Int("guild_count", len(batch)),
		slog.Int("total_roles", len(allRoles)),
		slog.Int("total_channels", len(allChannels)),
		slog.Int("total_emojis", len(allEmojis)),
		slog.Int("total_stickers", len(allStickers)),
	)

	return nil
}

func (b *GuildCreateBatcher) Flush(ctx context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.flushLocked(ctx)
}

func (b *GuildCreateBatcher) Close() error {
	b.cancel()
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.flushTimer != nil {
		b.flushTimer.Stop()
	}
	return b.flushLocked(context.Background())
}
