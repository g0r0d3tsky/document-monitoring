package config

import (
	"fmt"
	"github.com/caarlos0/env/v9"
)

type Config struct {
	MongoDB struct {
		Host     string `env:"MONGODB_HOST"`
		Port     int    `env:"MONGODB_PORT"`
		UserName string `env:"MONGODB_USER"`
		Password string `env:"MONGODB_PASSWORD"`
		DBName   string `env:"MONGODB_NAME"`
	}
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

func (c *Config) MongoDBDSN() string {
	hostPort := fmt.Sprintf("%s:%d", c.MongoDB.Host, c.MongoDB.Port)

	if c.MongoDB.UserName != "" && c.MongoDB.Password != "" {
		auth := fmt.Sprintf("%s:%s@", c.MongoDB.UserName, c.MongoDB.Password)
		return fmt.Sprintf("mongodb://%s%s", auth, hostPort)
	}

	return fmt.Sprintf("mongodb://%s", hostPort)
}
func Read() (*Config, error) {
	var config Config

	if err := env.Parse(&config); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return &config, nil
}

func (c *Config) PostgresDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Postgres.Host, c.Postgres.Port, c.Postgres.User, c.Postgres.Password, c.Postgres.Database,
	)
}
