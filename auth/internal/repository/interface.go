package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"time"
)

type DB interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type Repository interface {
	UserExists(ctx context.Context, username, email string) (bool, error)
	CreateUser(ctx context.Context, username, email, hash, role string) error
	GetUser(ctx context.Context, username string, hashedPassword *string, userId *int, role *string) error
	SaveRefreshToken(ctx context.Context, id int, token string, at time.Time) error
	GetRefreshTokenInfo(ctx context.Context, token string) (int, time.Time, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	GetUserForGenerateNewToken(ctx context.Context, userId int, username, email, role *string) error
}
