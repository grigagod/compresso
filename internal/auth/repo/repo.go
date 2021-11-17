package repo

import (
	"context"

	"github.com/grigagod/compresso/internal/auth"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type AuthRepo struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) *AuthRepo {
	return &AuthRepo{db}
}

// InsertUser insert new user in DB.
func (r *AuthRepo) InsertUser(ctx context.Context, user *auth.User) (*auth.User, error) {
	query := `INSERT INTO svc.users(user_id, username, password) VALUES($1, $2, $3) RETURNING *`
	u := &auth.User{}

	if err := r.db.QueryRowxContext(ctx, query, &user.ID, &user.Username, &user.Password).StructScan(u); err != nil {
		return nil, errors.Wrap(err, "AuthRepo.InsertUser.StructScan")
	}
	return u, nil
}

// SelectUserByName find user by name in DB.
func (r *AuthRepo) SelectUserByName(ctx context.Context, username string) (*auth.User, error) {
	query := `SELECT * FROM svc.users WHERE username = $1`
	u := &auth.User{}

	if err := r.db.QueryRowxContext(ctx, query, username).StructScan(u); err != nil {
		return nil, errors.Wrap(err, "AuthRepo.GetUserByName.StructScan")
	}

	return u, nil
}
