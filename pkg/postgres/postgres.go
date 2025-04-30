package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string `yaml:"POSTGRES_HOST" env:"POSTGRES_HOST" env-default:"localhost"`
	Port     uint16 `yaml:"POSTGRES_PORT" env:"POSTGRES_PORT" env-default:"5432"`
	Username string `yaml:"POSTGRES_USER" env:"POSTGRES_USER" env-default:"root"`
	Password string `yaml:"POSTGRES_PASSWORD" env:"POSTGRES_PASSWORD" env-default:"1234"`
	Database string `yaml:"POSTGRES_DB" env:"POSTGRES_DB" env-default:"postgres"`
}

func New(ctx context.Context, config Config) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)
	fmt.Println(connString)
	conn, err := pgxpool.New(ctx, connString)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}
	m, err := migrate.New("file://db/migrations", connString)
	if err != nil {
		return nil, fmt.Errorf("unable to create migration instance: %w", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}
	return conn, nil
}
