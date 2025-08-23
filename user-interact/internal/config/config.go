package config

import (
	"eventify/common/kafka"
	"fmt"

	"eventify/common/postgres"
	"github.com/ilyakaznacheev/cleanenv"
)

type UserInteractConfig struct {
	Port          int    `yaml:"port"`
	MigrationPath string `yaml:"migration-path"`
	Secret        string `yaml:"secret-key"`
}

type Config struct {
	Postgres     postgres.Config    `yaml:"postgres"`
	UserInteract UserInteractConfig `yaml:"user-interact"`
	Kafka        kafka.Config       `yaml:"kafka"`
}

func New() (Config, error) {
	var config Config
	// docker workdir app/
	// local workdir delivery-tracker/auth
	if err := cleanenv.ReadConfig("../configs/config.yaml", &config); err != nil {
		fmt.Println(err)
		if err := cleanenv.ReadEnv(&config); err != nil {
			return Config{}, fmt.Errorf("error reading configs: %w", err)
		}
	}

	return config, nil
}
