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
	InstanceCount int `toml:"instance_count" validate:"required"`
	InstanceIndex int `toml:"instance_index"`
}
