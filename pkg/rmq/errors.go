package rmq

import "errors"

// Pkg defined errors.
var (
	ErrNotFoundMethod     = errors.New("not found method")
	ErrWithRequeue        = errors.New("error with requeue")
	ErrNotEnoughArguments = errors.New("error not enough arguments")
	ErrWrongValue         = errors.New("error wrong value")
)
