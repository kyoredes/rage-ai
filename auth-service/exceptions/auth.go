package exceptions

import "errors"

var (
	ErrRefreshTokenNotFound   = errors.New("Refresh token not found")
	ErrAccessTokenNotFound    = errors.New("Access token not found")
	ErrAccessTokenNotCreated  = errors.New("Access token not created")
	ErrRefreshTokenNotCreated = errors.New("Refresh token not created")
	ErrDeleteRefreshToken     = errors.New("Error deleting refresh token")
)
