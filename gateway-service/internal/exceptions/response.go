package exceptions

import "errors"

var (
	ErrResponseExternalService = errors.New("Error response from external service")
	ErrUserNotFound            = errors.New("User not found")
)
