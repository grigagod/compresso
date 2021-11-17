//go:build unit

package repo

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/grigagod/compresso/internal/auth"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestAuthRepo_InsertUser(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authRepo := NewAuthRepository(sqlxDB)

	t.Run("Main case", func(t *testing.T) {
		user := &auth.User{
			ID:        uuid.New(),
			Username:  "test",
			Password:  "test",
			CreatedAt: time.Now(),
		}

		rows := sqlmock.NewRows([]string{"user_id", "username", "password", "created_at"}).AddRow(
			&user.ID, "test", "test", &user.CreatedAt)

		query := `INSERT INTO svc.users(user_id, username, password) VALUES($1, $2, $3) RETURNING *`
		mock.ExpectQuery(query).
			WithArgs(&user.ID, &user.Username, &user.Password).
			WillReturnRows(rows)

		createdUser, err := authRepo.InsertUser(context.Background(), user)

		require.NoError(t, err)
		require.Equal(t, createdUser, user)

	})
}

func TestAuthRepo_SelectUserByName(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authRepo := NewAuthRepository(sqlxDB)

	t.Run("Main Case", func(t *testing.T) {
		user := &auth.User{
			ID:        uuid.New(),
			Username:  "test",
			Password:  "test",
			CreatedAt: time.Now(),
		}

		rows := sqlmock.NewRows([]string{"user_id", "username", "password", "created_at"}).AddRow(
			&user.ID, "test", "test", &user.CreatedAt)

		query := `SELECT * FROM svc.users WHERE username = $1`
		mock.ExpectQuery(query).WithArgs(user.Username).WillReturnRows(rows)

		foundUser, err := authRepo.SelectUserByName(context.Background(), user.Username)

		require.NoError(t, err)
		require.Equal(t, foundUser, user)
	})
}
