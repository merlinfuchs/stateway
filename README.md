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
