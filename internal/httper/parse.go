package httper

import (
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

func ParseValidatorError(err error) Error {
	if strings.Contains(err.Error(), "Password") {
		return NewStatusMsg(http.StatusBadRequest, InvalidPasswordMsg)
	}

	if strings.Contains(err.Error(), "Username") {
		return NewStatusMsg(http.StatusBadRequest, InvalidUsernameMsg)
	}

	return NewStatusError(http.StatusBadRequest, err)
}

func ParseSqlError(err error) Error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			return NewBadRequestMsg(UserExistsMsg)
		case pgerrcode.NoDataFound:
			return NewBadRequestMsg(UserNotFoundMsg)
		}
	}

	return NewBadRequestError(err)
}

func ParseJWTError(err error) Error {
	var jwtErr *jwt.ValidationError
	if errors.As(err, &jwtErr) {
		switch jwtErr.Errors {
		case jwt.ValidationErrorExpired:
			return NewStatusMsg(http.StatusUnauthorized, TokenExpiredMsg)
		}
	}

	return NewStatusMsg(http.StatusUnauthorized, InvalidTokenMsg)
}
