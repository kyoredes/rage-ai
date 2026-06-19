package grpcserver

import (
	"auth/internal/logging"
	"auth/internal/service"
	"context"

	authv1 "agrobot/proto/gen/go/auth/v1"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	authv1.UnimplementedAuthServiceServer
	authService   *service.AuthService
	adminService  *service.AdminService
	healthService *service.HealthService
}

func NewAuthServer(
	authService *service.AuthService,
	adminService *service.AdminService,
	healthService *service.HealthService,
) *AuthServer {
	return &AuthServer{
		authService:   authService,
		adminService:  adminService,
		healthService: healthService,
	}
}

func (s *AuthServer) StartTelegram(ctx context.Context, req *authv1.StartTelegramRequest) (*authv1.StartTelegramResponse, error) {
	if req.GetTelegramId() == "" {
		return nil, status.Error(codes.InvalidArgument, "telegram_id is required")
	}

	result, err := s.authService.StartTelegramWithUser(req.GetTelegramId())
	if err != nil {
		logging.Logger.Error("grpc StartTelegram failed", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to start telegram user")
	}

	return &authv1.StartTelegramResponse{
		UserId:       result.UserID,
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}, nil
}

func (s *AuthServer) GetTelegramProfile(ctx context.Context, req *authv1.GetTelegramProfileRequest) (*authv1.GetTelegramProfileResponse, error) {
	if req.GetTelegramId() == "" {
		return nil, status.Error(codes.InvalidArgument, "telegram_id is required")
	}

	profile, err := s.authService.GetTelegramProfile(req.GetTelegramId())
	if err != nil {
		logging.Logger.Error("grpc GetTelegramProfile failed", zap.Error(err))
		return nil, status.Error(codes.NotFound, "telegram profile not found")
	}

	return &authv1.GetTelegramProfileResponse{
		UserId:     profile.UserID,
		TelegramId: profile.TelegramID,
		Email:      profile.Email,
	}, nil
}

func (s *AuthServer) ListUsers(ctx context.Context, req *authv1.ListUsersRequest) (*authv1.ListUsersResponse, error) {
	users, total, err := s.adminService.ListUsers(int(req.GetPage()), int(req.GetLimit()), req.GetSearch())
	if err != nil {
		logging.Logger.Error("grpc ListUsers failed", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to list users")
	}

	items := make([]*authv1.UserItem, len(users))
	for i, u := range users {
		items[i] = &authv1.UserItem{
			UserId:     u.UserID,
			Email:      u.Email,
			TelegramId: u.TelegramID,
			CreatedAt:  u.CreatedAt.Unix(),
		}
	}

	return &authv1.ListUsersResponse{
		Users: items,
		Total: int32(total),
	}, nil
}

func (s *AuthServer) GetUser(ctx context.Context, req *authv1.GetUserRequest) (*authv1.GetUserResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	user, err := s.adminService.GetUser(req.GetUserId())
	if err != nil {
		logging.Logger.Error("grpc GetUser failed", zap.Error(err))
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &authv1.GetUserResponse{
		UserId:     user.UserID,
		Email:      user.Email,
		TelegramId: user.TelegramID,
		CreatedAt:  user.CreatedAt.Unix(),
		UpdatedAt:  user.UpdatedAt.Unix(),
	}, nil
}

func (s *AuthServer) UpdateUser(ctx context.Context, req *authv1.UpdateUserRequest) (*authv1.UpdateUserResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	user, err := s.adminService.UpdateUser(req.GetUserId(), req.GetEmail())
	if err != nil {
		logging.Logger.Error("grpc UpdateUser failed", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to update user")
	}

	return &authv1.UpdateUserResponse{
		UserId: user.UserID,
		Email:  user.Email,
	}, nil
}

func (s *AuthServer) DeleteUser(ctx context.Context, req *authv1.DeleteUserRequest) (*authv1.DeleteUserResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	if err := s.adminService.DeleteUser(req.GetUserId()); err != nil {
		logging.Logger.Error("grpc DeleteUser failed", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete user")
	}

	return &authv1.DeleteUserResponse{}, nil
}

func (s *AuthServer) GetStats(ctx context.Context, req *authv1.GetAuthStatsRequest) (*authv1.GetAuthStatsResponse, error) {
	stats, err := s.adminService.GetStats()
	if err != nil {
		logging.Logger.Error("grpc GetStats failed", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get stats")
	}

	return &authv1.GetAuthStatsResponse{
		TotalUsers: int32(stats.TotalUsers),
		NewUsers_7D: int32(stats.NewUsers7d),
	}, nil
}

func (s *AuthServer) Health(ctx context.Context, req *authv1.HealthRequest) (*authv1.HealthResponse, error) {
	dbOk, redisOk := s.healthService.Check(ctx)
	return &authv1.HealthResponse{
		Ok:      dbOk && redisOk,
		DbOk:    dbOk,
		RedisOk: redisOk,
	}, nil
}
