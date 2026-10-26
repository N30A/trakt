package auth

import "errors"

var (
	ErrInvalidToken       = errors.New("invalid token")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInitialUserExists  = errors.New("initial user already exists")
)
