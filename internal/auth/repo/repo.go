package repo

import (
	"github.com/grigagod/compresso/internal/auth"
	"github.com/grigagod/compresso/internal/auth/models"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type authRepo struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) auth.Repository {
	return &authRepo{db}
}

func (r *authRepo) Create(user *models.User) (*models.User, error) {
	u := &models.User{}

	if err := r.db.QueryRowx(createUserQuery, &user.ID, &user.Username, &user.Password, &user.CreatedAt).StructScan(u); err != nil {
		return nil, errors.Wrap(err, "authRepo.Create.StructScan")
	}
	return u, nil
}

func (r *authRepo) FindByName(username string) (*models.User, error) {
	u := &models.User{}

	if err := r.db.QueryRowx(findUserByNameQuery, username).StructScan(u); err != nil {
		return nil, errors.Wrap(err, "authRepo.FindByName.StructScan")
	}

	return u, nil
}
