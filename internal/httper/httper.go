package httper

import (
	"errors"
	"net/http"
)

type Error interface {
	error
	Status() int
}

type StatusError struct {
	Code int
	Err  error
}

func NewStatusError(code int, err error) StatusError {
	return StatusError{
		Code: code,
		Err:  err,
	}
}

func NewStatusMsg(code int, msg ErrMessage) StatusError {
	return StatusError{
		Code: code,
		Err:  errors.New(string(msg)),
	}
}

func NewBadRequestError(err error) StatusError {
	return NewStatusError(http.StatusBadRequest, err)
}

func NewBadRequestMsg(msg ErrMessage) StatusError {
	return NewStatusMsg(http.StatusBadRequest, msg)
}

func NewUnauthorizedError(err error) StatusError {
	return NewStatusError(http.StatusUnauthorized, err)
}

func NewUnauthorizedMsg(msg ErrMessage) StatusError {
	return NewStatusMsg(http.StatusUnauthorized, msg)
}

func NewWrongCredentialsMsg() StatusError {
	return NewUnauthorizedMsg(WrongCredentialsMsg)
}

func NewNotFoundMsg() StatusError {
	return NewStatusMsg(http.StatusNotFound, NotFoundMsg)
}

func NewNotAllowedMediaMsg() StatusError {
	return NewStatusMsg(http.StatusUnsupportedMediaType, NotAllowedMediaTypeMsg)
}

func (st StatusError) Error() string {
	return st.Err.Error()
}

func (st StatusError) Status() int {
	return st.Code
}
