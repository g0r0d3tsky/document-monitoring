package config

import (
	"fmt"
	"github.com/caarlos0/env/v9"
)

type Config struct {
	Host     string `env:"MONGODB_HOST"`
	Port     int    `env:"MONGODB_PORT"`
	UserName string `env:"MONGODB_USER"`
	Password string `env:"MONGODB_PASSWORD"`
	DBName   string `env:"MONGODB_NAME"`
}

func (c *Config) MongoDBDSN() string {
	hostPort := fmt.Sprintf("%s:%d", c.Host, c.Port)

	if c.UserName != "" && c.Password != "" {
		auth := fmt.Sprintf("%s:%s@", c.UserName, c.Password)
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
