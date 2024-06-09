package db

import (
	"goqueue/pkg/xerrors"
)

var (
	ErrCritical = xerrors.New("critical error")
	ErrNotFound = xerrors.New("not found")
)
