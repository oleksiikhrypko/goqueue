package db

import (
	"errors"

	"goqueue/pkg/xerrors"
)

var (
	ErrCritical = xerrors.New(errors.New("critical error"))
	ErrNotFound = xerrors.New(errors.New("not found"))
)
