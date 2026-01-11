package app

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/disgoorg/snowflake/v2"
	"github.com/merlinfuchs/stateway/stateway-gateway/model"
	"github.com/merlinfuchs/stateway/stateway-gateway/store"
	"github.com/merlinfuchs/stateway/stateway-lib/event"
)

type AppManagerConfig struct {
	GatewayCount int
	GatewayID    int
	NoResume     bool
}

type AppManager struct {
	sync.Mutex

	cfg                    AppManagerConfig
	appStore               store.AppStore
	groupStore             store.GroupStore
	shardSessionStore      store.ShardSessionStore
	identifyRateLimitStore store.IdentifyRateLimitStore
	eventHandler           event.EventHandler

	groups map[string]*model.Group
	apps   map[snowflake.ID]*App
}

func NewAppManager(
	cfg AppManagerConfig,
	appStore store.AppStore,
	groupStore store.GroupStore,
	shardSessionStore store.ShardSessionStore,
	identifyRateLimitStore store.IdentifyRateLimitStore,
	eventHandler event.EventHandler,
) *AppManager {
	// Discord some times sends unquoted snowflake IDs, so we need to allow them
	snowflake.AllowUnquoted = true

	return &AppManager{
		cfg:                    cfg,
		appStore:               appStore,
		groupStore:             groupStore,
		shardSessionStore:      shardSessionStore,
		identifyRateLimitStore: identifyRateLimitStore,
		eventHandler:           eventHandler,

		apps:   make(map[snowflake.ID]*App),
		groups: make(map[string]*model.Group),
	}
}

func (m *AppManager) Run(ctx context.Context) {
	m.populateGroups(ctx)
	m.populateApps(ctx, time.Time{})

	lastUpdate := time.Now()
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(10 * time.Second):
			m.populateGroups(ctx)
			m.populateApps(ctx, lastUpdate)
			lastUpdate = time.Now()
		}
	}
}

func (m *AppManager) populateGroups(ctx context.Context) {
	groups, err := m.groupStore.GetGroups(ctx)
	if err != nil {
		slog.Error("Failed to get groups", slog.Any("error", err))
		return
	}

	m.Lock()
	defer m.Unlock()

	for _, group := range groups {
		m.groups[group.ID] = group
	}
}

func (m *AppManager) populateApps(ctx context.Context, lastUpdate time.Time) {
	apps, err := m.appStore.GetEnabledApps(ctx, store.GetEnabledAppsParams{
		GatewayCount: m.cfg.GatewayCount,
		GatewayID:    m.cfg.GatewayID,
	})
	if err != nil {
		slog.Error("Failed to get enabled apps", slog.Any("error", err))
		return
	}

	for _, app := range apps {
		if app.UpdatedAt.After(lastUpdate) {
			m.addOrUpdateApp(ctx, app)
		}
	}

	m.removeDanglingApps(ctx, apps)
}

func (m *AppManager) addOrUpdateApp(ctx context.Context, app *model.App) {
	m.Lock()
	defer m.Unlock()

	group, ok := m.groups[app.GroupID]
	if !ok {
		slog.Error("Group not found", slog.String("group_id", app.GroupID))
		return
	}

	if _, ok := m.apps[app.ID]; ok {
		m.apps[app.ID].Update(ctx, app, group)
	} else {
		newApp := NewApp(
			AppConfig(m.cfg),
			app,
			group,
			m.appStore,
			m.shardSessionStore,
			m.identifyRateLimitStore,
			m.eventHandler,
		)
		m.apps[app.ID] = newApp
		go newApp.Run(ctx)
	}
}

func (m *AppManager) removeDanglingApps(ctx context.Context, apps []*model.App) {
	appIDs := make(map[snowflake.ID]bool)
	for _, app := range apps {
		appIDs[app.ID] = true
	}

	m.Lock()
	defer m.Unlock()

	for id, app := range m.apps {
		if !appIDs[id] {
			slog.Info("Removing dangling app", slog.String("app_id", id.String()))
			app.Close(ctx)
			delete(m.apps, id)
		}
	}
}
