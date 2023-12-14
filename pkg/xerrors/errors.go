package xerrors

import (
	"errors"
)

type Error struct {
	err        error
	message    string
	fieldsErrs map[string]string
}

func New(err error) *Error {
	return &Error{
		err: err,
	}
}

func (e *Error) Error() string {
	if e.message != "" {
		return e.message
	}
	if e.err == nil {
		return "operation failed"
	}
	return e.err.Error()
}

func (e *Error) Unwrap() error {
	return e.err
}

func (e *Error) Consume(err error) *Error {
	var (
		message    string
		fieldsErrs map[string]string
	)
	if e != nil {
		message = e.message
		fieldsErrs = e.fieldsErrs
	}
	var xerr *Error
	if errors.As(err, &xerr) {
		message = xerr.message
		fieldsErrs = xerr.fieldsErrs
	}

	return &Error{
		err:        errors.Join(err, e),
		message:    message,
		fieldsErrs: fieldsErrs,
	}
}

func (e *Error) WithMessage(message string) *Error {
	return &Error{
		err:     e,
		message: message,
	}
}

func (e *Error) WithDetails(fieldsErrs map[string]string) *Error {
	return &Error{
		err:        e,
		message:    e.message,
		fieldsErrs: fieldsErrs,
	}
}

func (e *Error) WithAdditionalInfo(message string, fieldsErrs map[string]string) *Error {
	return &Error{
		err:        e,
		message:    message,
		fieldsErrs: fieldsErrs,
	}
}

func (e *Error) FieldsErrors() map[string]string {
	if e == nil {
		return map[string]string{}
	}
	return e.fieldsErrs
}

func (e *Error) Message() string {
	if e == nil {
		return ""
	}
	return e.message
}
