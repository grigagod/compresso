//go:generate mockgen -source usecase.go -destination mock/usecase_mock.go -package mock
package auth

type UseCase interface {
	Register(user *User) (*UserWithToken, error)
	Login(user *User) (*UserWithToken, error)
}
