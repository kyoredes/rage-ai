package exceptions

import "errors"

var (
	ErrUserNotFound       = errors.New("User not found")
	ErrInvalidCredentials = errors.New("Invalid credentials")
	ErrUserAlreadyExists  = errors.New("User already exists")
	ErrCreatingUser       = errors.New("Error creating user")
)
