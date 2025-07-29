package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"vizurth/eventify/common/postgres"
)

type UserInteractConfig struct {
	Port          int    `yaml:"port" env:"PORT"`
	MigrationPath string `yaml:"migration-path" env:"MIGRATION_PATH"`
	Secret        string `yaml:"secret" env:"SECRET"`
}

type Config struct {
	Postgres     postgres.Config    `yaml:"postgres" env-prefix:"POSTGRES_"`
	UserInteract UserInteractConfig `yaml:"user-interact" env-previx:"EVENT_"`
}

func New() (Config, error) {
	var config Config
	// docker workdir app/
	// local workdir eventify/user-interact
	if err := cleanenv.ReadConfig("configs/configs.yaml", &config); err != nil {
		fmt.Println(err)
		if err := cleanenv.ReadEnv(&config); err != nil {
			return Config{}, fmt.Errorf("error reading configs: %w", err)
		}
	}

	return config, nil
}
