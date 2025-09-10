package repository

import (
	"context"
	"errors"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRepository_UserExists(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	un, em := "username", "email"

	rows := pgxmock.NewRows([]string{"count"}).AddRow(1)
	mock.ExpectQuery(`SELECT count\(\*\) FROM users WHERE email = \$1 AND username = \$2`).
		WithArgs(em, un).
		WillReturnRows(rows)

	repo := NewAuthRepository(mock)

	flag, err := repo.UserExists(context.Background(), un, em)

	require.NoError(t, err)
	require.True(t, flag)
}

func TestRepository_UserExists_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()
	un, em := "username", "email"

	mock.ExpectQuery(`SELECT count\(\*\) FROM users WHERE email = \$1 AND username = \$2`).
		WithArgs(em, un).
		WillReturnError(errors.New("error"))
	repo := NewAuthRepository(mock)
	flag, err := repo.UserExists(context.Background(), un, em)
	require.Error(t, err)
	require.False(t, flag)
}

func TestRepository_CreateUser(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	un, em := "username", "email"
	password := "password"
	role := "user"
	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO users \(username, email, password, role\) VALUES \(\$1, \$2, \$3, \$4\)`).
		WithArgs(un, em, password, role).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectCommit()
	repo := NewAuthRepository(mock)
	err = repo.CreateUser(context.Background(), un, em, password, role)

	require.NoError(t, err)
}
