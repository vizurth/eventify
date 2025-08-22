package repository

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	db   *pgxpool.Pool
	psql sq.StatementBuilderType
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{
		db:   db,
		psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *AuthRepository) UserExists(ctx context.Context, username, email string) (bool, error) {
	var count int

	query, args, err := r.psql.Select("count(*)").
		From("users").Where(sq.Eq{"email": email, "username": username}).ToSql()

	err = r.db.QueryRow(ctx, query, args...).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *AuthRepository) CreateUser(ctx context.Context, username, email, hash, role string) error {
	query, args, err := r.psql.Insert("users").Columns("username", "email", "password_hash", "role").Values(username, email, hash, role).ToSql()
	_, err = r.db.Exec(ctx, query, args...)

	if err != nil {
		return errors.New("could not create user")
	}
	return nil
}

func (r *AuthRepository) GetUser(ctx context.Context, username string, hashedPassword *string, userId *int, role *string) error {
	query, args, err := r.psql.Select("id", "password_hash", "role").From("users").Where(sq.Eq{"username": username}).ToSql()

	// Выполняем запрос
	row := r.db.QueryRow(ctx, query, args...)

	// Сканируем результат в переданные указатели
	err = row.Scan(userId, hashedPassword, role)
	if err != nil {
		return errors.New("could not get user")
	}

	return nil
}
