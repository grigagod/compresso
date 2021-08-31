package httper

type ErrMessage string

const (
	BadRequestMsg       ErrMessage = "Bad request"
	UserExistsMsg       ErrMessage = "User with such username already exists"
	UserNotFoundMsg     ErrMessage = "User with such username is not found"
	NotAllowedHeader    ErrMessage = "Not allowed header"
	WrongCredentialsMsg ErrMessage = "Wrong credentials"
	InvalidUsernameMsg  ErrMessage = "Username allowed length is 4-30 characters"
	InvalidPasswordMsg  ErrMessage = "Password allowed length is 4-40 characters"
	TokenExpiredMsg     ErrMessage = "Your token is expired, login again"
)
