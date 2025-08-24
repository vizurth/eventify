package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type GatewayConfig struct {
	Port            int           `yaml:"port" env:"PORT" env-default:"9090"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env:"SHUTDOWN_TIMEOUT" env-default:"30s"`
}

type AuthConfig struct {
	SecretKey string `yaml:"secret_key" env:"SECRET_KEY" env-default:"VxCDGAeugvEa6jvBiFYcOboddRGKydns"`
	URL       string `yaml:"url" env:"URL" env-default:"http://auth-service:9091"`
}

type EventConfig struct {
	URL string `yaml:"url" env:"URL" env-default:"http://event-service:9092"`
}

type UserInteractionConfig struct {
	URL string `yaml:"url" env:"URL" env-default:"http://user-interact-service:9093"`
}

type NotificationConfig struct {
	URL string `yaml:"url" env:"URL" env-default:"http://notification-service:9094"`
}

type Config struct {
	Gateway         GatewayConfig         `yaml:"gateway"`
	Auth            AuthConfig            `yaml:"auth"`
	Event           EventConfig           `yaml:"event"`
	UserInteraction UserInteractionConfig `yaml:"user_interact"`
	Notification    NotificationConfig    `yaml:"notification"`
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
