package client

import (
	"fmt"
	"gateway/internal/config"

	aiv1 "agrobot/proto/gen/go/ai/v1"
	authv1 "agrobot/proto/gen/go/auth/v1"
	subscriptionv1 "agrobot/proto/gen/go/subscription/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Clients struct {
	Auth         authv1.AuthServiceClient
	Subscription subscriptionv1.SubscriptionServiceClient
	AI           aiv1.AIServiceClient
	authConn     *grpc.ClientConn
	subConn      *grpc.ClientConn
	aiConn       *grpc.ClientConn
}

func NewClients(authCfg *config.AuthConfig, subCfg *config.SubConfig, aiCfg *config.AIConfig) (*Clients, error) {
	authConn, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", authCfg.AuthHost, authCfg.AuthGRPCPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("auth grpc dial: %w", err)
	}

	subConn, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", subCfg.SubHost, subCfg.SubGRPCPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		authConn.Close()
		return nil, fmt.Errorf("subscription grpc dial: %w", err)
	}

	aiConn, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", aiCfg.AIHost, aiCfg.AIGRPCPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		authConn.Close()
		subConn.Close()
		return nil, fmt.Errorf("ai grpc dial: %w", err)
	}

	return &Clients{
		Auth:         authv1.NewAuthServiceClient(authConn),
		Subscription: subscriptionv1.NewSubscriptionServiceClient(subConn),
		AI:           aiv1.NewAIServiceClient(aiConn),
		authConn:     authConn,
		subConn:      subConn,
		aiConn:       aiConn,
	}, nil
}

func (c *Clients) Close() error {
	if err := c.authConn.Close(); err != nil {
		return err
	}
	if err := c.subConn.Close(); err != nil {
		return err
	}
	return c.aiConn.Close()
}
