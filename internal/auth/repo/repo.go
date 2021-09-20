package repo

import (
	"github.com/grigagod/compresso/internal/auth"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type authRepo struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) *authRepo {
	return &authRepo{db}
}

func (r *authRepo) Create(user *auth.User) (*auth.User, error) {
	query := `INSERT INTO svc.users(user_id, username, password, created_at) VALUES($1, $2, $3, $4) RETURNING *`
	u := &auth.User{}

	if err := r.db.QueryRowx(query, &user.ID, &user.Username, &user.Password, &user.CreatedAt).StructScan(u); err != nil {
		return nil, errors.Wrap(err, "authRepo.Create.StructScan")
	}
	return u, nil
}

func (r *authRepo) FindByName(username string) (*auth.User, error) {
	query := `SELECT * FROM svc.users WHERE username = $1`
	u := &auth.User{}

	if err := r.db.QueryRowx(query, username).StructScan(u); err != nil {
		return nil, errors.Wrap(err, "authRepo.FindByName.StructScan")
	}

	return u, nil
}
