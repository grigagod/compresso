package httper

import "net/http"

type Error interface {
	error
	Status() int
}

type StatusError struct {
	Code int
	Err  string
}

func NewStatusError(code int, err string) StatusError {
	return StatusError{
		Code: code,
		Err:  err,
	}
}

func NewStatusMsg(code int, msg ErrMessage) StatusError {
	return StatusError{
		Code: code,
		Err:  string(msg),
	}
}

func NewBadRequestMsg(msg ErrMessage) StatusError {
	return StatusError{
		Code: http.StatusBadRequest,
		Err:  string(msg),
	}
}

func (st StatusError) Error() string {
	return st.Err
}

func (st StatusError) Status() int {
	return st.Code
}
