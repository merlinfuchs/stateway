package config

import (
	"github.com/go-playground/validator/v10"
)

type GatewayConfig struct {
	Logging  LoggingConfig  `toml:"logging"`
	Database DatabaseConfig `toml:"database"`
	Broker   BrokerConfig   `toml:"broker"`
}

func (cfg *GatewayConfig) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate.Struct(cfg)
}

type CacheConfig struct {
	Logging LoggingConfig `toml:"logging"`
	Broker  BrokerConfig  `toml:"broker"`
}

func (cfg *CacheConfig) Validate() error {
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
