package klist

import "github.com/pkg/errors"

var (
	ErrCritical = errors.New("klist: critical db error")
	ErrNotFound = errors.New("klist: not found")
)
