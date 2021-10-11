//go:generate mockgen -source usecase.go -destination mock/usecase_mock.go -package mock
package auth

import "context"

type UseCase interface {
	Register(ctx context.Context, user *User) (*UserWithToken, error)
	Login(ctx context.Context, user *User) (*UserWithToken, error)
}
