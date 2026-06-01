package core

import "errors"

var (
	ErrNotFound   = errors.New("not found")
	ErrBadCreds   = errors.New("bad credentials")
	ErrBadToken   = errors.New("bad token")
	ErrUserExists = errors.New("username taken")
)