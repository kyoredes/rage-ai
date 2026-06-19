package grpcserver

import (
	"auth/internal/logging"
	"auth/internal/service"
	"context"

	authv1 "rageai/proto/gen/go/auth/v1"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	authv1.UnimplementedAuthServiceServer
	authService *service.AuthService
}

func NewAuthServer(authService *service.AuthService) *AuthServer {
	return &AuthServer{authService: authService}
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
