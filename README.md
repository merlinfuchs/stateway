# Stateway

Stateway is a toolkit for building scalable Discord bots.

## Services

### Gateway

Connects all (or a subset of) the configured applications to the Discord gateway and dispatches events over NATS.

It primarily uses the `GATEWAY` stream and `gateway.>` subjects.

### Cache

Receives events from the Gateway over NATS and stores entities like guilds, roles, channels, etc. in a PostgreSQL database.

It responds to requests on the `cache.>` subjects.

## Library

The `stateway-lib` package contains the core libraries for Stateway. It can be used by clients to interact with the Stateway services.
