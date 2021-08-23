package repo

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/grigagod/compresso/internal/auth/models"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestAuthRepo_Create(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authRepo := NewAuthRepository(sqlxDB)

	t.Run("Create", func(t *testing.T) {
		user := &models.User{
			ID:        uuid.New(),
			Username:  "test",
			Password:  "test",
			CreatedAt: time.Now(),
		}

		rows := sqlmock.NewRows([]string{"user_id", "username", "password", "created_at"}).AddRow(
			&user.ID, "test", "test", &user.CreatedAt)

		mock.ExpectQuery(createUserQuery).
			WithArgs(&user.ID, &user.Username, &user.Password, &user.CreatedAt).
			WillReturnRows(rows)

		createdUser, err := authRepo.Create(user)

		require.NoError(t, err)
		require.Equal(t, createdUser, user)

	})
}

func TestAuthRepo_FindByName(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authRepo := NewAuthRepository(sqlxDB)

	t.Run("FindByName", func(t *testing.T) {
		user := &models.User{
			ID:        uuid.New(),
			Username:  "test",
			Password:  "test",
			CreatedAt: time.Now(),
		}

		rows := sqlmock.NewRows([]string{"user_id", "username", "password", "created_at"}).AddRow(
			&user.ID, "test", "test", &user.CreatedAt)

		mock.ExpectQuery(findUserByNameQuery).WithArgs(user.Username).WillReturnRows(rows)

		foundUser, err := authRepo.FindByName(user.Username)

		require.NoError(t, err)
		require.Equal(t, foundUser, user)
	})
}
