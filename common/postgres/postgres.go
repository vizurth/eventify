package postgres

import (
	"context"
	"errors"
	"eventify/common/logger"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
	"time"
)

type Config struct {
	Host     string `yaml:"host" env:"POSTGRES_HOST" env-default:"0.0.0.0"`
	Port     uint16 `yaml:"port" env:"POSTGRES_PORT" env-default:"5432"`
	Username string `yaml:"username" env:"POSTGRES_USER" env-default:"root"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD" env-default:"1234"`
	Database string `yaml:"database" env:"POSTGRES_DB" env-default:"postgres"`
}

//type Config struct {
//	Host     string `yaml:"host" env:"POSTGRES_HOST" `
//	Port     uint16 `yaml:"port" env:"POSTGRES_PORT" `
//	Username string `yaml:"username" env:"POSTGRES_USER" `
//	Password string `yaml:"password" env:"POSTGRES_PASSWORD"`
//	Database string `yaml:"database" env:"POSTGRES_DB"`
//}

func New(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	connString := cfg.GetConnString()

	conn, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	return conn, nil
}
func WaitForPostgres(ctx context.Context, cfg Config, retries int, delay time.Duration) error {
	connString := cfg.GetConnString()
	var err error
	for i := 0; i < retries; i++ {
		var conn *pgxpool.Pool
		conn, err = pgxpool.New(ctx, connString)
		if err == nil {
			conn.Close()
			return nil
		}
		logger.GetLoggerFromCtx(ctx).Info(ctx, "waiting for postgres to be ready...")
		time.Sleep(delay)
	}
	return fmt.Errorf("postgres is not ready after %d retries: %w", retries, err)
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

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		// Если ошибка указывает на dirty базу
		if strings.Contains(err.Error(), "Dirty database") {
			version, dirty, versionErr := m.Version()
			if versionErr != nil {
				return fmt.Errorf("migration dirty error, but failed to get version: %w", versionErr)
			}
			if dirty {
				if forceErr := m.Force(int(version)); forceErr != nil {
					return fmt.Errorf("failed to force migration version: %w", forceErr)
				}
				// Повторяем Up() после очистки dirty-флага
				err = m.Up()
				if err != nil && !errors.Is(err, migrate.ErrNoChange) {
					return fmt.Errorf("failed to run migrations after force: %w", err)
				}
			}
		} else {
			return fmt.Errorf("failed to run migrations: %w", err)
		}
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
