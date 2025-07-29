package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) UserExists(ctx context.Context, username, email string) (bool, error) {
	var count int

	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM schema_name.users WHERE username = $1 OR email = $2", username, email).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *AuthRepository) CreateUser(ctx context.Context, username, email, hash, role string) error {
	_, err := r.db.Exec(ctx, "INSERT INTO schema_name.users (username, email, password_hash, role) VALUES ($1, $2, $3, $4)", username, email, hash, role)

	if err != nil {
		return errors.New("could not create user")
	}
	return nil
}

func (r *AuthRepository) GetUser(ctx context.Context, username string, hashedPassword *string, userId *int, role *string) error {
	query := `
		SELECT id, password_hash, role
		FROM schema_name.users
		WHERE username = $1 OR email = $1
		LIMIT 1
	`

	// Выполняем запрос
	row := r.db.QueryRow(ctx, query, username)

	// Сканируем результат в переданные указатели
	err := row.Scan(userId, hashedPassword, role)
	if err != nil {
		return errors.New("could not get user")
	}

	return nil
}
