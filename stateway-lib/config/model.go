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

type RootAuditConfig struct {
	Logging  LoggingConfig  `toml:"logging"`
	Database DatabaseConfig `toml:"database"`
	Broker   BrokerConfig   `toml:"broker"`
	Audit    AuditConfig    `toml:"audit"`
}

func (cfg *RootAuditConfig) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate.Struct(cfg)
}

type AuditConfig struct {
	GatewayIDs []int `toml:"gateway_ids"`
}

type DatabaseConfig struct {
	Postgres   PostgresConfig   `toml:"postgres"`
	Clickhouse ClickhouseConfig `toml:"clickhouse"`
}

type LoggingConfig struct {
	Filename   string `toml:"filename"`
	MaxSize    int    `toml:"max_size"`
	MaxAge     int    `toml:"max_age"`
	MaxBackups int    `toml:"max_backups"`
	Debug      bool   `toml:"debug"`
}

type PostgresConfig struct {
	Host     string `toml:"host" validate:"required"`
	Port     int    `toml:"port" validate:"required"`
	DBName   string `toml:"db_name" validate:"required"`
	User     string `toml:"user" validate:"required"`
	Password string `toml:"password"`
}

type ClickhouseConfig struct {
	Host     string `toml:"host" validate:"required"`
	Port     int    `toml:"port" validate:"required"`
	DBName   string `toml:"db_name" validate:"required"`
	User     string `toml:"user" validate:"required"`
	Password string `toml:"password"`
}

type BrokerConfig struct {
	NATS          NATSConfig `toml:"nats"`
	NamePrefix    string     `toml:"name_prefix"`
	SubjectPrefix string     `toml:"subject_prefix"`
}

type NATSConfig struct {
	URL string `toml:"url" validate:"required"`
}

type GatewayConfig struct {
	GatewayCount int                  `toml:"gateway_count" validate:"required"`
	GatewayID    int                  `toml:"gateway_id"`
	Groups       []GatewayGroupConfig `toml:"groups"`
	Apps         []GatewayAppConfig   `toml:"apps"`
	NoResume     bool                 `toml:"no_resume"`
}

type GatewayGroupConfig struct {
	ID          string `toml:"id" validate:"required"`
	DisplayName string `toml:"display_name" validate:"required"`
	MaxShards   int    `toml:"max_shards"`
	MaxGuilds   int    `toml:"max_guilds"`
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
	InMemory   bool  `toml:"in_memory"`
	GatewayIDs []int `toml:"gateway_ids"`
}
