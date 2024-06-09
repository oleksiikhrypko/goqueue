package xerrors

import (
	"errors"
)

type Error struct {
	err        error
	message    string
	extensions map[string]any
}

func New(msg string) *Error {
	return &Error{
		message: "",
		err:     errors.New(msg),
	}
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.message != "" {
		return e.message
	}
	if e.err != nil {
		return e.err.Error()
	}
	return "action failed"
}

func (e *Error) Unwrap() error {
	return e.err
}

func (e *Error) Consume(err error) *Error {
	var (
		message string
	)
	if e != nil {
		message = e.message
	}
	var xerr *Error
	if errors.As(err, &xerr) {
		message = xerr.message
	}

	return &Error{
		err:     errors.Join(err, e),
		message: message,
	}
}

func (e *Error) WithMessage(message string) *Error {
	if e == nil {
		return &Error{
			message: message,
		}
	}
	return &Error{
		err:     e,
		message: message,
	}
}

func (e *Error) WithExtensions(extensions map[string]any) *Error {
	ext := make(map[string]any)
	for k, v := range extensions {
		ext[k] = v
	}
	if e == nil {
		return &Error{
			extensions: ext,
		}
	}
	return &Error{
		err:        e,
		message:    e.message,
		extensions: ext,
	}
}

func (e *Error) WithAdditionalInfo(message string, extensions map[string]any) *Error {
	ext := make(map[string]any)
	for k, v := range extensions {
		ext[k] = v
	}
	if e == nil {
		return &Error{
			message:    message,
			extensions: ext,
		}
	}
	return &Error{
		err:        e,
		message:    message,
		extensions: ext,
	}
}

func (e *Error) Extensions() map[string]any {
	if e == nil {
		return map[string]any{}
	}

	ext := make(map[string]any)
	for k, v := range e.extensions {
		ext[k] = v
	}

	err := errors.Unwrap(e)
	switch x := err.(type) {
	case interface{ Unwrap() error }:
		extractExtensions(x.Unwrap(), ext)
	case interface{ Unwrap() []error }:
		for _, err := range x.Unwrap() {
			extractExtensions(err, ext)
		}
	}
	return ext
}

func extractExtensions(err error, dest map[string]any) {
	var xerr *Error
	if errors.As(err, &xerr) {
		for k, v := range xerr.Extensions() {
			if _, ok := dest[k]; ok {
				continue
			}
			dest[k] = v
		}
	}
}

func (e *Error) Message() string {
	if e == nil {
		return ""
	}
	return e.message
}
