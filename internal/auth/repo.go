//go:generate mockgen -source repo.go -destination mock/repo_mock.go -package mock
package auth

import "context"

type Repository interface {
	Create(ctx context.Context, user *User) (*User, error)
	FindByName(ctx context.Context, username string) (*User, error)
}
