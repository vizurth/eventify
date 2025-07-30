package config

import (
	"fmt"
	"os"

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
}

func New() (Config, error) {
	var config Config

	const configPath = "configs/config.yaml"

	// Проверим, существует ли файл

	if _, err := os.Stat(configPath); err != nil {
		fmt.Println("config file not found, falling back to env:", err)
		if err := cleanenv.ReadEnv(&config); err != nil {
			return Config{}, fmt.Errorf("error reading config from env: %w", err)
		}
	} else {
		if err := cleanenv.ReadConfig(configPath, &config); err != nil {
			return Config{}, fmt.Errorf("error reading config from file: %w", err)
		}
	}

	// Логируем результат
	fmt.Printf("✅ LOADED CONFIG: %+v\n", config)

	return config, nil
}
