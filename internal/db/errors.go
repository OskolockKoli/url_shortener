package db

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrDuplicateKey = errors.New("duplicate key")
)
