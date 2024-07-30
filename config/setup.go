package config

import (
	"context"

	handler "github.com/ADAGroupTcc/ms-channels-api/internal/http/channels"
	"github.com/ADAGroupTcc/ms-channels-api/internal/http/health"
	repository "github.com/ADAGroupTcc/ms-channels-api/internal/repositories/channels"
	service "github.com/ADAGroupTcc/ms-channels-api/internal/services/channels"
	healthService "github.com/ADAGroupTcc/ms-channels-api/internal/services/health"
	"github.com/ADAGroupTcc/ms-channels-api/pkg/mongorm"
)

type Dependencies struct {
	Handler       handler.Handler
	HealthHandler health.Health
}

func NewDependencies(ctx context.Context, envs *Environments) *Dependencies {
	database, err := mongorm.Connect(envs.DBUri, envs.DBName)
	if err != nil {
		panic(err)
	}
	channelsRepository := repository.New(database)
	channelService := service.New(channelsRepository)
	channelHandler := handler.New(channelService)

	healthService := healthService.New(database)
	healthHandler := health.New(healthService)
	return &Dependencies{
		channelHandler,
		healthHandler,
	}
}
