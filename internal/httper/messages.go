package httper

type ErrMessage string

const (
	BadRequestMsg          ErrMessage = "Bad request"
	UserExistsMsg          ErrMessage = "User with such username already exists"
	NotFoundMsg            ErrMessage = "Not found"
	NotAllowedHeaderMsg    ErrMessage = "Provoded header is not allowed"
	NotAllowedContentType  ErrMessage = "Provided content type is not allowed"
	WrongCredentialsMsg    ErrMessage = "Wrong credentials"
	InvalidUsernameMsg     ErrMessage = "Username allowed length is 4-30 characters"
	InvalidPasswordMsg     ErrMessage = "Password allowed length is 4-40 characters"
	UnexpectedSignatureMsg ErrMessage = "Unexpected signature method"
	InvalidTokenMsg        ErrMessage = "Provided token is invalid"
	TokenExpiredMsg        ErrMessage = "Provided token is expired, login again"
)
