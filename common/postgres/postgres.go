package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string `yaml:"host" env:"POSTGRES_HOST"`
	Port     uint16 `yaml:"port" env:"POSTGRES_PORT"`
	Username string `yaml:"username" env:"POSTGRES_USER"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD"`
	Database string `yaml:"database" env:"POSTGRES_DB"`
}

func New(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	connString := cfg.GetConnString()

	conn, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	return conn, nil
}

func (c *Config) GetConnString() string {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
	)
	return connString
}
