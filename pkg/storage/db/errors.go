package db

import "errors"

var (
	ErrCritical = errors.New("critical db error")
	ErrNotFound = errors.New("not found")
)
