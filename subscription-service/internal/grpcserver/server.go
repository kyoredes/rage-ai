package grpcserver

import (
	"context"
	"subscription/internal/logging"
	"subscription/internal/service"
	"time"

	subscriptionv1 "agrobot/proto/gen/go/subscription/v1"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SubscriptionServer struct {
	subscriptionv1.UnimplementedSubscriptionServiceServer
	service       *service.SubscriptionService
	adminService  *service.AdminService
	healthService *service.HealthService
}

func NewSubscriptionServer(
	service *service.SubscriptionService,
	adminService *service.AdminService,
	healthService *service.HealthService,
) *SubscriptionServer {
	return &SubscriptionServer{
		service:       service,
		adminService:  adminService,
		healthService: healthService,
	}
}

func (s *SubscriptionServer) GetSubscriptionByUserId(ctx context.Context, req *subscriptionv1.GetSubscriptionByUserIdRequest) (*subscriptionv1.GetSubscriptionByUserIdResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	sub, err := s.service.EnsureSubByUserId(userID)
	if err != nil {
		logging.Logger.Error("grpc GetSubscriptionByUserId failed", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get subscription")
	}

	return &subscriptionv1.GetSubscriptionByUserIdResponse{
		SubscriptionId: sub.Uuid.String(),
		UserId:         sub.UserID.String(),
		StartsAt:       sub.StartsAt.Unix(),
		ExpiresAt:      sub.ExpiresAt.Unix(),
	}, nil
}

func (s *SubscriptionServer) ListSubscriptions(ctx context.Context, req *subscriptionv1.ListSubscriptionsRequest) (*subscriptionv1.ListSubscriptionsResponse, error) {
	subs, total, err := s.adminService.ListSubscriptions(int(req.GetPage()), int(req.GetLimit()), req.GetStatus())
	if err != nil {
		logging.Logger.Error("grpc ListSubscriptions failed", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to list subscriptions")
	}

	items := make([]*subscriptionv1.SubscriptionItem, len(subs))
	for i, sub := range subs {
		items[i] = &subscriptionv1.SubscriptionItem{
			SubscriptionId: sub.Uuid.String(),
			UserId:         sub.UserID.String(),
			StartsAt:       sub.StartsAt.Unix(),
			ExpiresAt:      sub.ExpiresAt.Unix(),
		}
	}

	return &subscriptionv1.ListSubscriptionsResponse{
		Subscriptions: items,
		Total:         int32(total),
	}, nil
}

func (s *SubscriptionServer) UpdateSubscription(ctx context.Context, req *subscriptionv1.UpdateSubscriptionRequest) (*subscriptionv1.UpdateSubscriptionResponse, error) {
	if req.GetSubscriptionId() == "" {
		return nil, status.Error(codes.InvalidArgument, "subscription_id is required")
	}

	startsAt := time.Unix(req.GetStartsAt(), 0)
	expiresAt := time.Unix(req.GetExpiresAt(), 0)

	sub, err := s.adminService.UpdateSubscription(req.GetSubscriptionId(), startsAt, expiresAt)
	if err != nil {
		logging.Logger.Error("grpc UpdateSubscription failed", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &subscriptionv1.UpdateSubscriptionResponse{
		SubscriptionId: sub.Uuid.String(),
		UserId:         sub.UserID.String(),
		StartsAt:       sub.StartsAt.Unix(),
		ExpiresAt:      sub.ExpiresAt.Unix(),
	}, nil
}

func (s *SubscriptionServer) DeleteSubscription(ctx context.Context, req *subscriptionv1.DeleteSubscriptionRequest) (*subscriptionv1.DeleteSubscriptionResponse, error) {
	if req.GetSubscriptionId() == "" {
		return nil, status.Error(codes.InvalidArgument, "subscription_id is required")
	}

	if err := s.adminService.DeleteSubscription(req.GetSubscriptionId()); err != nil {
		logging.Logger.Error("grpc DeleteSubscription failed", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete subscription")
	}

	return &subscriptionv1.DeleteSubscriptionResponse{}, nil
}

func (s *SubscriptionServer) GetStats(ctx context.Context, req *subscriptionv1.GetSubscriptionStatsRequest) (*subscriptionv1.GetSubscriptionStatsResponse, error) {
	stats, err := s.adminService.GetStats()
	if err != nil {
		logging.Logger.Error("grpc GetStats failed", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get stats")
	}

	return &subscriptionv1.GetSubscriptionStatsResponse{
		Total:   int32(stats.Total),
		Active:  int32(stats.Active),
		Expired: int32(stats.Expired),
	}, nil
}

func (s *SubscriptionServer) Health(ctx context.Context, req *subscriptionv1.HealthRequest) (*subscriptionv1.HealthResponse, error) {
	dbOk := s.healthService.Check(ctx)
	return &subscriptionv1.HealthResponse{
		Ok:   dbOk,
		DbOk: dbOk,
	}, nil
}
