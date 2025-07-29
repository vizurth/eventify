package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"vizurth/eventify/common/logger"
)

type Config struct {
	Host     string `yaml:"host" env:"POSTGRES_HOST" env-default:"0.0.0.0"`
	Port     uint16 `yaml:"port" env:"POSTGRES_PORT" env-default:"5432"`
	Username string `yaml:"username" env:"POSTGRES_USER" env-default:"root"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD" env-default:"1234"`
	Database string `yaml:"database" env:"POSTGRES_DB" env-default:"postgres"`
}

func New(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	connString := cfg.GetConnString()

	conn, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	return conn, nil
}

func Migrate(ctx context.Context, cfg Config, migrationsPath string) error {
	connString := cfg.GetConnString()

	m, err := migrate.New(
		migrationsPath,
		connString,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize migrations: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	logger.GetLoggerFromCtx(ctx).Info(ctx, "migrated successfully")
	return nil
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
