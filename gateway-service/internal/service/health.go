package service

import (
	"context"
	"gateway/internal/dto"
	"gateway/internal/logging"
	"time"

	authv1 "agrobot/proto/gen/go/auth/v1"
	aiv1 "agrobot/proto/gen/go/ai/v1"
	subscriptionv1 "agrobot/proto/gen/go/subscription/v1"

	"go.uber.org/zap"
)

func (s *AdminService) GetServicesStatus() (*dto.ServicesStatusResponse, error) {
	authStatus, authDB, authRedis, authLatency := s.checkAuthHealth()
	subStatus, subDB, subLatency := s.checkSubscriptionHealth()
	aiStatus, aiRedis, aiLatency := s.checkAIHealth()

	redisUp := authRedis || aiRedis

	services := []dto.ServiceStatus{
		{ID: "gateway", Name: "Gateway", Status: "up", LatencyMs: 0},
		{ID: "auth-service", Name: "Auth Service", Status: authStatus, LatencyMs: authLatency},
		{ID: "auth-db", Name: "Auth DB", Status: boolStatus(authDB), LatencyMs: authLatency},
		{ID: "subscription-service", Name: "Subscription Service", Status: subStatus, LatencyMs: subLatency},
		{ID: "sub-db", Name: "Subscription DB", Status: boolStatus(subDB), LatencyMs: subLatency},
		{ID: "ai-service", Name: "AI Service", Status: aiStatus, LatencyMs: aiLatency},
		{ID: "redis", Name: "Redis", Status: boolStatus(redisUp), LatencyMs: redisLatency(authLatency, aiLatency)},
	}

	return &dto.ServicesStatusResponse{
		Services:  services,
		CheckedAt: time.Now().Unix(),
	}, nil
}

func boolStatus(ok bool) string {
	if ok {
		return "up"
	}
	return "down"
}

func redisLatency(authMs, aiMs int64) int64 {
	if authMs == 0 {
		return aiMs
	}
	if aiMs == 0 {
		return authMs
	}
	if authMs < aiMs {
		return authMs
	}
	return aiMs
}

func (s *AdminService) checkAuthHealth() (status string, dbOk, redisOk bool, latencyMs int64) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	start := time.Now()
	resp, err := s.clients.Auth.Health(ctx, &authv1.HealthRequest{})
	latencyMs = time.Since(start).Milliseconds()
	if err != nil {
		logging.Logger.Error("auth Health failed", zap.Error(err))
		return "down", false, false, latencyMs
	}

	dbOk = resp.GetDbOk()
	redisOk = resp.GetRedisOk()
	if resp.GetOk() {
		return "up", dbOk, redisOk, latencyMs
	}
	return "degraded", dbOk, redisOk, latencyMs
}

func (s *AdminService) checkSubscriptionHealth() (status string, dbOk bool, latencyMs int64) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	start := time.Now()
	resp, err := s.clients.Subscription.Health(ctx, &subscriptionv1.HealthRequest{})
	latencyMs = time.Since(start).Milliseconds()
	if err != nil {
		logging.Logger.Error("subscription Health failed", zap.Error(err))
		return "down", false, latencyMs
	}

	dbOk = resp.GetDbOk()
	if resp.GetOk() {
		return "up", dbOk, latencyMs
	}
	return "degraded", dbOk, latencyMs
}

func (s *AdminService) checkAIHealth() (status string, redisOk bool, latencyMs int64) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	start := time.Now()
	resp, err := s.clients.AI.Health(ctx, &aiv1.HealthRequest{})
	latencyMs = time.Since(start).Milliseconds()
	if err != nil {
		logging.Logger.Error("ai Health failed", zap.Error(err))
		return "down", false, latencyMs
	}

	redisOk = resp.GetRedisOk()
	if resp.GetOk() {
		return "up", redisOk, latencyMs
	}
	return "degraded", redisOk, latencyMs
}
