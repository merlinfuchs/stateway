# Stateway (WIP)

Stateway is a toolkit for building scalable Discord bots. This is heavily work in progress and not yet ready to be used outside of [Embed Generator](https://github.com/merlinfuchs/embed-generator).

## Services

### Gateway

Connects all (or a subset of) the configured applications to the Discord gateway and dispatches events over NATS.

It primarily uses the `GATEWAY` stream and `gateway.>` subjects.

### Cache

Receives events from the Gateway over NATS and stores entities like guilds, roles, channels, etc. in a PostgreSQL database.

It responds to requests on the `cache.>` subjects.

## Library

The `stateway-lib` package contains the core libraries for Stateway. It can be used by clients to interact with the Stateway services.

## NATS Infra

### GATEWAY Stream

The `GATEWAY` receives and stores events from the `stateway-gateway` service. It is primarily used to forward Discord gateway events to any other service that needs to know about them.

In the future it may also contain `stateway-gateway` specific custom events.

#### Subject Structure

The subject is composed of the following parts:

`gateway.<gateway_id>.<group_id>.<app_id>.<...event_type>`

The event type is the lower case version of the Discord event type with underscores replaced by dots. This makes it possible to match on event type groups. (e.g. `GUILD_CREATE` -> `guild.create`)

Example: `gateway.0.default.1234567890.guild.create`

Example matches:

- `gateway.0.default.1234567890.guild.create` only matches if the gateway_id is 0, the group_id is `default`, the app_id is `1234567890` and the event_type is `guild.create`
- `gateway.*.*.*.guild.create` matches for `guild.create` events from any gateway, group or app
- `gateway.*.*.*.guild.>` matches for `guild.create`, `guild.update`, `guild.delete`, etc. events from any gateway, group or app
- `gateway.0.>` matches for all events from gateway 0
- `gateway.>` matches for all events from any gateway, group or app

## Configuration

All Stateway services will read their configuration from a `stateway.toml` file in the current working directory.

Some service specific configuration are scoped under a `[<server_type>]` section.

```toml
[logging]
filename = "stateway.log" # The filename of the log file to write to. Leave empty to only log to stdout.
max_size = 100 # The maximum size of the log file in megabytes.
max_age = 7 # The maximum age of the log file in days.
max_backups = 10 # The maximum number of backup log files to keep.

[broker.nats]
url = "nats://127.0.0.1:4222" # The URL of the NATS server to connect to.

[database.postgres]
host = "127.0.0.1" # The host of the PostgreSQL server to connect to.
port = 5432 # The port of the PostgreSQL server to connect to.
user = "postgres" # The user to connect to the PostgreSQL server with.
db_name = "stateway" # The database to connect to.

# Stateway Gateway configuration.
[gateway]
gateway_count = 1 # The number of gateways you are running and balance the apps across.
gateway_id = 0 # The ID of the gateway you are running (0-based index).

# Any apps you want to always run, apps can also be dynamically added and removed using the admin CLI.
[[gateway.apps]]
token = "your-discord-bot-token" # Your Discord bot token.
shard_count = 2 # The number of shards to run for this app.
shard_concurrency = 1 # The maximum number of concurrent identify requests to send to the Discord gateway.
group_id = "default" # The group to run the app in.
intents = 1023 # The intents to use for this app.
# Optional presence configuration.
presence = {
    status = "online",
    activity = {
        name = "Stateway",
        state = "Running",
        type = "WATCHING",
        url = "https://github.com/merlinfuchs/stateway"
    }
}

[cache]
gateway_ids = [0, 1, 2] # The gateway IDs to process events from. Leave empty to process events from all gateways.
```
