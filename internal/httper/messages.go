package httper

type ErrMessage string

const (
	BadRequestMsg          ErrMessage = "Bad request"
	UserExistsMsg          ErrMessage = "User with such username already exists"
	UserNotFoundMsg        ErrMessage = "User with such username is not found"
	NotAllowedHeaderMsg    ErrMessage = "Provoded header is not allowed"
	NotAllowedMediaTypeMsg ErrMessage = "Provided media type is not allowed"
	NotAllowedFileSizeMsg  ErrMessage = "Provided file size is not allowed"
	WrongCredentialsMsg    ErrMessage = "Wrong credentials"
	InvalidUsernameMsg     ErrMessage = "Username allowed length is 4-30 characters"
	InvalidPasswordMsg     ErrMessage = "Password allowed length is 4-40 characters"
	UnexpectedSignatureMsg ErrMessage = "Unexpected signature method"
	InvalidTokenMsg        ErrMessage = "Provided token is invalid"
	TokenExpiredMsg        ErrMessage = "Provided token is expired, login again"
)
