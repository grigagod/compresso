package httper

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"
)

func ParseValidatorError(err error) Error {
	if strings.Contains(err.Error(), "Password") {
		return NewStatusMsg(http.StatusBadRequest, InvalidPasswordMsg)
	}

	if strings.Contains(err.Error(), "Username") {
		return NewStatusMsg(http.StatusBadRequest, InvalidUsernameMsg)
	}

	return NewStatusError(http.StatusBadRequest, err.Error())
}

func ParseSqlError(err error) Error {
	if errors.Is(err, sql.ErrNoRows) {
		return NewStatusMsg(http.StatusNotFound, UserNotFoundMsg)
	}
	if strings.Contains(err.Error(), "23505") {
		return NewStatusMsg(http.StatusBadRequest, UserExistsMsg)
	}

	return NewStatusError(http.StatusBadRequest, err.Error())
}
