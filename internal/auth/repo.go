//go:generate mockgen -source repo.go -destination mock/repo_mock.go -package mock
package auth

import "context"

type Repository interface {
	InsertUser(ctx context.Context, user *User) (*User, error)
	SelectUserByName(ctx context.Context, username string) (*User, error)
}
