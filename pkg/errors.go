package pkg

import "errors"

var (
	ErrNoRows         = errors.New("no rows in result set")
	ErrNoRowsAffected = errors.New("no rows affected")
	ErrDuplicate      = errors.New("duplicate key")
)
