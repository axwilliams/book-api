package web

import (
	"errors"
)

var (
	ErrInternalServer = errors.New("Something went wrong, we are aware of the problem")
	ErrValidation     = errors.New("Validation error")
)

type Validation struct {
	Message string   `json:"message"`
	Errors  []string `json:"errors"`
}

type RequestError struct {
	Err    error
	Status int
	Fields []string
}

func NewRequestError(err error, status int) error {
	return &RequestError{err, status, nil}
}

func (re *RequestError) Error() string {
	return re.Err.Error()
}
