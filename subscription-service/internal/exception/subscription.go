package exception

import "errors"

var (
	ErrSubscriptionNotFound       = errors.New("subscription not found")
	ErrSubscriptionAlreadyExists  = errors.New("subscription already exists")
	ErrCreatingSubscription       = errors.New("error creating subscription")
	ErrInvalidSubscriptionDates   = errors.New("expires_at must be after starts_at")
)
