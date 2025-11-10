package config

import (
	"github.com/go-playground/validator/v10"
)

type RootGatewayConfig struct {
	Logging  LoggingConfig  `toml:"logging"`
	Database DatabaseConfig `toml:"database"`
	Broker   BrokerConfig   `toml:"broker"`
	Gateway  GatewayConfig  `toml:"gateway"`
}

func (cfg *RootGatewayConfig) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate.Struct(cfg)
}

type RootCacheConfig struct {
	Logging  LoggingConfig  `toml:"logging"`
	Database DatabaseConfig `toml:"database"`
	Broker   BrokerConfig   `toml:"broker"`
	Cache    CacheConfig    `toml:"cache"`
}

func (cfg *RootCacheConfig) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate.Struct(cfg)
}

type DatabaseConfig struct {
	Postgres PostgresConfig `toml:"postgres"`
}

type LoggingConfig struct {
	Filename   string `toml:"filename"`
	MaxSize    int    `toml:"max_size"`
	MaxAge     int    `toml:"max_age"`
	MaxBackups int    `toml:"max_backups"`
}

type PostgresConfig struct {
	Host     string `toml:"host" validate:"required"`
	Port     int    `toml:"port" validate:"required"`
	DBName   string `toml:"db_name" validate:"required"`
	User     string `toml:"user" validate:"required"`
	Password string `toml:"password"`
}

type BrokerConfig struct {
	NATS NATSConfig `toml:"nats"`
}

type NATSConfig struct {
	URL string `toml:"url" validate:"required"`
}

type GatewayConfig struct {
	GatewayCount int                `toml:"gateway_count" validate:"required"`
	GatewayID    int                `toml:"gateway_id"`
	Apps         []GatewayAppConfig `toml:"apps"`
}

type GatewayAppConfig struct {
	Token            string                    `toml:"token" validate:"required"`
	ShardCount       int                       `toml:"shard_count" validate:"required,min=1"`
	ShardConcurrency int                       `toml:"shard_concurrency"`
	GroupID          string                    `toml:"group_id"`
	Intents          int64                     `toml:"intents"`
	Presence         *GatewayAppPresenceConfig `toml:"presence"`
}

type GatewayAppPresenceConfig struct {
	Status   string                            `toml:"status"`
	Activity *GatewayAppPresenceActivityConfig `toml:"activity"`
}

type GatewayAppPresenceActivityConfig struct {
	Name  string `toml:"name"`
	State string `toml:"state"`
	Type  string `toml:"type"`
	URL   string `toml:"url"`
}

type CacheConfig struct {
	GatewayIDs []int `toml:"gateway_ids"`
}
