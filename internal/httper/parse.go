package httper

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/grigagod/compresso/pkg/converter"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

func ParseValidatorError(err error) Error {
	if strings.Contains(err.Error(), "Password") {
		return NewBadRequestMsg(InvalidPasswordMsg)
	}

	if strings.Contains(err.Error(), "Username") {
		return NewBadRequestMsg(InvalidUsernameMsg)
	}

	if strings.Contains(err.Error(), "CRF") {
		return NewBadRequestError(converter.ErrVideoCRF)
	}

	return NewBadRequestError(err)
}

func ParseSqlError(err error) Error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.UniqueViolation {
			return NewBadRequestMsg(UserExistsMsg)
		}
	}

	if errors.Is(err, sql.ErrNoRows) {
		return NewStatusMsg(http.StatusNotFound, NotFoundMsg)
	}

	return NewBadRequestError(err)
}

func ParseJWTError(err error) Error {
	var jwtErr *jwt.ValidationError
	if errors.As(err, &jwtErr) {
		if jwtErr.Errors == jwt.ValidationErrorExpired {
			return NewStatusMsg(http.StatusUnauthorized, TokenExpiredMsg)
		}
	}

	return NewUnauthorizedMsg(InvalidTokenMsg)
}
