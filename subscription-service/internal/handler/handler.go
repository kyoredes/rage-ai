package handler

import "subscription/internal/service"

type Handler struct {
	Subscription *SubHandler
}

func NewHandler(
	subscriptionService *service.SubscriptionService,
) *Handler {
	return &Handler{
		Subscription: NewSubHandler(subscriptionService),
	}
}
