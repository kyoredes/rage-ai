package exception

import "errors"

var (
	ErrSubscriptionNotFound      = errors.New("subscription not found")
	ErrSubscriptionAlreadyExists = errors.New("subscription already exists")
	ErrCreatingSubscription      = errors.New("error creating subscription")
)
