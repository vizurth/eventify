package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server       ServerConfig              `yaml:"server" env-prefix:"SERVER_"`
	Auth         AuthServiceConfig         `yaml:"auth" env-prefix:"AUTH_"`
	Event        EventServiceConfig        `yaml:"event" env-prefix:"EVENT_"`
	UserInteract UserInteractServiceConfig `yaml:"user-interact" env-prefix:"USER_INTERACT_"`
}

type ServerConfig struct {
	Port int `yaml:"port" env:"PORT" env-default:"8080"`
}

type AuthServiceConfig struct {
	Host string `yaml:"host" env:"HOST" env-default:"localhost"`
	Port int    `yaml:"port" env:"PORT" env-default:"5001"`
}

type EventServiceConfig struct {
	Host string `yaml:"host" env:"HOST" env-default:"localhost"`
	Port int    `yaml:"port" env:"PORT" env-default:"5002"`
}

type UserInteractServiceConfig struct {
	Host string `yaml:"host" env:"HOST" env-default:"localhost"`
	Port int    `yaml:"port" env:"PORT" env-default:"5003"`
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
