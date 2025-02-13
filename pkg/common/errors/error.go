package errors

import (
	"errors"
	"fmt"
)

type BadRequestError struct {
	Wrapper
}

type Wrapper struct {
	Base  error
	Msg   string
	Debug any
}

func (e *Wrapper) Error() string {
	if e.Msg != "" {
		return e.Msg
	}
	if e.Base != nil {
		return e.Base.Error()
	}
	return ""
}

func (e *Wrapper) WithDebug(debug any) *Wrapper {
	e.Debug = debug
	return e
}

func New(template string, args ...any) error {
	return fmt.Errorf(template, args...)
}

func BadRequest(template string, args ...any) *BadRequestError {
	return &BadRequestError{
		Wrapper{Msg: fmt.Sprintf(template, args...)},
	}
}

func Is(err error, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}
