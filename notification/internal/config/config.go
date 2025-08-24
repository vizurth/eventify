package config

import (
	"eventify/common/kafka"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type NotificationConfig struct {
	Port int `yaml:"port" env:"PORT" env-default:"9095"`
}

type Config struct {
	Kafka        kafka.Config       `yaml:"kafka"`
	Notification NotificationConfig `yaml:"notification"`
}

func New() (Config, error) {
	var config Config
	if err := cleanenv.ReadConfig("../configs/config.yaml", &config); err != nil {
		fmt.Println(err)
		if err := cleanenv.ReadEnv(&config); err != nil {
			return Config{}, fmt.Errorf("error reading configs: %w", err)
		}
	}
	return config, nil
}
