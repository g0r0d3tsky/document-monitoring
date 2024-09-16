package config

import (
	"fmt"
	"github.com/caarlos0/env/v9"
)

type Config struct {
	Postgres struct {
		Host     string `env:"POSTGRES_HOST,notEmpty"`
		Port     string `env:"POSTGRES_PORT,notEmpty"`
		User     string `env:"POSTGRES_USER,notEmpty"`
		Password string `env:"POSTGRES_PASSWORD,notEmpty"`
		Database string `env:"POSTGRES_DB,notEmpty"`
	}
	Kafka struct {
		BrokerList string `env:"KAFKA_BROKERS,notEmpty"`
		GroupID    string `env:"KAFKA_GROUP_ID,notEmpty"`
		Topics     string `env:"KAFKA_TOPICS,notEmpty"`
	}
}

func (c *Config) PostgresDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Postgres.Host, c.Postgres.Port, c.Postgres.User, c.Postgres.Password, c.Postgres.Database,
	)
}

func Read() (*Config, error) {
	var config Config
	if err := env.Parse(&config); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return &config, nil
}
