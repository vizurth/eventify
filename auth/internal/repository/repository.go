package repository

import (
	"context"
	"eventify/common/logger"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type AuthRepository struct {
	db    DB
	psql  sq.StatementBuilderType
	redis *redis.Client
}

func NewAuthRepository(db DB, redis *redis.Client) Repository {
	return &AuthRepository{
		db:    db,
		psql:  sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		redis: redis,
	}
}

// UserExists checking if user exists in users table
func (r *AuthRepository) UserExists(ctx context.Context, username, email string) (bool, error) {
	var count int

	query, args, err := r.psql.Select("count(*)").
		From("users").Where(sq.Eq{"email": email, "username": username}).ToSql()

	err = r.db.QueryRow(ctx, query, args...).Scan(&count)

	if err != nil {
		return false, fmt.Errorf("user exists repository error: %w", err)
	}

	return count > 0, nil
}

// CreateUser create in user in users table
func (r *AuthRepository) CreateUser(ctx context.Context, username, email, hash, role string) error {
	query, args, err := r.psql.Insert("users").Columns("username", "email", "password_hash", "role").Values(username, email, hash, role).ToSql()
	_, err = r.db.Exec(ctx, query, args...)

	if err != nil {
		return fmt.Errorf("create user repository error: %w", err)
	}
	return nil
}

// GetUser get user from table
func (r *AuthRepository) GetUser(ctx context.Context, username string, hashedPassword *string, userId *int, role *string) error {
	query, args, err := r.psql.Select("id", "password_hash", "role").From("users").Where(sq.Eq{"username": username}).ToSql()

	// Выполняем запрос
	row := r.db.QueryRow(ctx, query, args...)

	// Сканируем результат в переданные указатели
	err = row.Scan(userId, hashedPassword, role)
	if err != nil {
		return fmt.Errorf("get user repository error: %w", err)
	}

	return nil
}

// GetUserForGenerateNewToken get user for generate token after expires
func (r *AuthRepository) GetUserForGenerateNewToken(ctx context.Context, userId int, username, email, role *string) error {
	query, args, err := r.psql.Select("username", "email", "role").
		From("users").
		Where(sq.Eq{"id": userId}).
		ToSql()
	if err != nil {
		return fmt.Errorf("get user repository error: %w", err)
	}
	err = r.db.QueryRow(ctx, query, args...).Scan(&username, email, role)
	if err != nil {
		return fmt.Errorf("get user repository error: %w", err)
	}

	return nil
}

// SaveRefreshToken save refresh token in table
func (r *AuthRepository) SaveRefreshToken(ctx context.Context, userId int, refreshToken string, expiresAt time.Time) error {
	// Postgres
	query, args, err := r.psql.Insert("refresh_tokens").
		Columns("user_id", "token", "expires_at").
		Values(userId, refreshToken, expiresAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("save refresh token repository error: %w", err)
	}
	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("save refresh token repository error: %w", err)
	}

	// Redis
	ttl := time.Until(expiresAt)
	err = r.redis.Set(ctx, fmt.Sprintf("refresh:%s", refreshToken), userId, ttl).Err()

	if err != nil {
		return fmt.Errorf("save refresh token repository error to redis: %w", err)
	}
	logger.GetOrCreateLoggerFromCtx(ctx).Info(ctx, "save refresh token to redis success")
	return nil
}

// GetRefreshTokenInfo get token from table
func (r *AuthRepository) GetRefreshTokenInfo(ctx context.Context, token string) (int, time.Time, error) {
	// Redis
	userIdStr, err := r.redis.Get(ctx, fmt.Sprintf("refresh:%s", token)).Result()
	if err == nil {
		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			return 0, time.Time{}, fmt.Errorf("invalid user id from redis: %w", err)
		}
		ttl, err := r.redis.TTL(ctx, fmt.Sprintf("refresh:%s", token)).Result()
		expiresAt := time.Now().Add(ttl)
		logger.GetOrCreateLoggerFromCtx(ctx).Info(ctx, "successfully get refresh token info from redis")
		return userId, expiresAt, nil
	}

	// Postgres
	var userId int
	var expiresAt time.Time
	query, args, err := r.psql.Select("user_id", "expires_at").
		From("refresh_tokens").
		Where(sq.Eq{"token": token}).ToSql()

	if err != nil {
		return 0, time.Time{}, fmt.Errorf("get refresh token repository error: %w", err)
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&userId, &expiresAt)
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("get refresh token repository error: %w", err)
	}
	return userId, expiresAt, nil
}

// DeleteRefreshToken delete token from table
func (r *AuthRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	// Postgres
	query, args, err := r.psql.Delete("refresh_tokens").Where(sq.Eq{"token": token}).ToSql()
	if err != nil {
		return fmt.Errorf("delete refresh token repository error: %w", err)
	}
	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("delete refresh token repository error: %w", err)
	}

	// Redis
	err = r.redis.Del(ctx, fmt.Sprintf("refresh:%s", token)).Err()
	if err != nil {
		return fmt.Errorf("delete refresh token repository error: %w", err)
	}

	return nil
}
