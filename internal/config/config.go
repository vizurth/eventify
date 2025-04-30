package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/vizurth/eventify/pkg/postgres"
)

type Config struct {
	SecretKey      string          `yaml:"SECRET_KEY" env:"SECRET_KEY" env-default:"your-secret-key"`
	PostgresConfig postgres.Config `yaml:"POSTGRES" env:"POSTGRES" env-default:"postgres"`
}

func NewConfig() (*Config, error) {
	var config Config
	if err := cleanenv.ReadConfig("./config/config.yaml", &config); err != nil {
		return nil, err
	}
	return &config, nil
}
