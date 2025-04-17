package users

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserParams   = errors.New("id and name cannot be empty")
)
