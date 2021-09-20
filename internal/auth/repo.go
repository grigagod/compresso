//go:generate mockgen -source repo.go -destination mock/repo_mock.go -package mock
package auth

type Repository interface {
	Create(user *User) (*User, error)
	FindByName(username string) (*User, error)
}
