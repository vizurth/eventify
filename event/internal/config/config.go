package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"vizurth/eventify/common/postgres"
)

type EventConfig struct {
	Port          int    `yaml:"port" env:"PORT"`
	MigrationPath string `yaml:"migration-path" env:"MIGRATION_PATH"`
}

type Config struct {
	Postgres postgres.Config `yaml:"postgres" env-prefix:"POSTGRES_"`
	Event    EventConfig     `yaml:"event" env-previx:"EVENT_"`
}

func New() (Config, error) {
	var config Config
	// docker workdir app/
	// local workdir delivery-tracker/event
	if err := cleanenv.ReadConfig("configs/configs.yaml", &config); err != nil {
		fmt.Println(err)
		if err := cleanenv.ReadEnv(&config); err != nil {
			return Config{}, fmt.Errorf("error reading configs: %w", err)
		}
	}

	return config, nil
}
