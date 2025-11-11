package server

import (
	"context"
	"errors"
	"time"

	"github.com/disgoorg/snowflake/v2"
	"github.com/merlinfuchs/stateway/stateway-gateway/store"
	"github.com/merlinfuchs/stateway/stateway-lib/gateway"
	"github.com/merlinfuchs/stateway/stateway-lib/service"
	"gopkg.in/guregu/null.v4"
)

type Gateway struct {
	groupStore store.GroupStore
	appStore   store.AppStore
}

func NewGateway(groupStore store.GroupStore, appStore store.AppStore) *Gateway {
	return &Gateway{
		groupStore: groupStore,
		appStore:   appStore,
	}
}

func (g *Gateway) GetApp(ctx context.Context, appID snowflake.ID) (*gateway.App, error) {
	app, err := g.appStore.GetApp(ctx, appID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, service.ErrNotFound("app not found")
		}
		return nil, err
	}
	return app, nil
}

func (g *Gateway) GetApps(ctx context.Context, params gateway.ListAppsRequest) ([]*gateway.App, error) {
	apps, err := g.appStore.GetApps(ctx, store.GetAppsParams{
		GroupID: params.GroupID,
		Limit:   params.Limit,
		Offset:  params.Offset,
	})
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, service.ErrNotFound("apps not found")
		}
		return nil, err
	}
	return apps, nil
}

func (g *Gateway) UpsertApp(ctx context.Context, params gateway.UpsertAppRequest) (*gateway.App, error) {
	app, err := g.appStore.UpsertApp(ctx, store.UpsertAppParams{
		ID:                  params.ID,
		GroupID:             params.GroupID,
		DisplayName:         params.DisplayName,
		DiscordClientID:     params.DiscordClientID,
		DiscordBotToken:     params.DiscordBotToken,
		DiscordPublicKey:    params.DiscordPublicKey,
		DiscordClientSecret: params.DiscordClientSecret,
		ShardCount:          params.ShardCount,
		Constraints:         params.Constraints,
		Config:              params.Config,
		CreatedAt:           time.Now().UTC(),
		UpdatedAt:           time.Now().UTC(),
	})
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (g *Gateway) DisableApp(ctx context.Context, appID snowflake.ID) error {
	err := g.appStore.DisableApp(ctx, store.DisableAppParams{
		ID:              appID,
		DisabledCode:    gateway.AppDisabledCodeUnknown,
		DisabledMessage: null.String{},
		UpdatedAt:       time.Now().UTC(),
	})
	if err != nil {
		return err
	}
	return nil
}

func (g *Gateway) DeleteApp(ctx context.Context, appID snowflake.ID) error {
	err := g.appStore.DeleteApp(ctx, appID)
	if err != nil {
		return err
	}
	return nil
}

func (g *Gateway) GetGroup(ctx context.Context, groupID string) (*gateway.Group, error) {
	group, err := g.groupStore.GetGroup(ctx, groupID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, service.ErrNotFound("group not found")
		}
		return nil, err
	}
	return group, nil
}

func (g *Gateway) GetGroups(ctx context.Context) ([]*gateway.Group, error) {
	groups, err := g.groupStore.GetGroups(ctx)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (g *Gateway) UpsertGroup(ctx context.Context, params gateway.UpsertGroupRequest) (*gateway.Group, error) {
	group, err := g.groupStore.UpsertGroup(ctx, store.UpsertGroupParams{
		ID:                 params.ID,
		DisplayName:        params.DisplayName,
		DefaultConstraints: params.DefaultConstraints,
		DefaultConfig:      params.DefaultConfig,
		CreatedAt:          time.Now().UTC(),
		UpdatedAt:          time.Now().UTC(),
	})
	if err != nil {
		return nil, err
	}
	return group, nil
}

func (g *Gateway) DeleteGroup(ctx context.Context, groupID string) error {
	err := g.groupStore.DeleteGroup(ctx, groupID)
	if err != nil {
		return err
	}
	return nil
}
