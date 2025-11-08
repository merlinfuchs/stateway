package app

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
	"github.com/merlinfuchs/stateway/stateway-gateway/model"
)

func shardsFromApp(app *model.App, gatewayCount int, gatewayID int) (int, []int) {
	shardCount := app.ShardCount
	if shardCount == 0 {
		shardCount = 1
	}

	shardIDs := make([]int, 0, shardCount)
	for shardID := 0; shardID < shardCount; shardID++ {
		// if shardCount == 1, this gateway is the only one that should run the app
		// otherwise, we are splitting the app shards across the gateways
		// shardID % gatewayCount gives us the index of the gateway that should run the shard
		if shardCount == 1 || shardID%gatewayCount == gatewayID {
			shardIDs = append(shardIDs, shardID)
		}
	}

	return shardCount, shardIDs
}

func intentsFromApp(app *model.App) gateway.Intents {
	intents := gateway.IntentsNonPrivileged
	if app.Config.Intents.Valid {
		intents = gateway.Intents(app.Config.Intents.Int64)
	}
	return intents
}

func presenceOptsFromApp(app *model.App) []gateway.PresenceOpt {
	res := []gateway.PresenceOpt{}
	if app.Config.Presence != nil {
		presence := app.Config.Presence
		if presence.Status.Valid {
			res = append(res, gateway.WithOnlineStatus(discord.OnlineStatus(presence.Status.String)))
		}

		if presence.Activity != nil {
			activityOpts := []gateway.ActivityOpt{}

			switch presence.Activity.Type {
			case "watching":
				res = append(
					res,
					gateway.WithWatchingActivity(presence.Activity.Name, activityOpts...),
				)
			case "listening":
				res = append(
					res,
					gateway.WithListeningActivity(presence.Activity.Name, activityOpts...),
				)
			case "competing":
				res = append(
					res,
					gateway.WithCompetingActivity(presence.Activity.Name, activityOpts...),
				)
			case "streaming":
				res = append(
					res,
					gateway.WithStreamingActivity(presence.Activity.Name, presence.Activity.URL, activityOpts...),
				)
			case "playing":
				res = append(
					res,
					gateway.WithPlayingActivity(presence.Activity.Name, activityOpts...),
				)
			default:
				res = append(
					res,
					gateway.WithCustomActivity(presence.Activity.Name, activityOpts...),
				)
			}
		}
	}

	return res
}
