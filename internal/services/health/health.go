package healthService

import (
	"context"

	"github.com/ADAGroupTcc/ms-channels-api/internal/services/health/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type HealthService interface {
	Check(ctx context.Context) domain.HealthResponse
}

type healthService struct {
	database *mongo.Database
}

func New(database *mongo.Database) HealthService {
	return &healthService{database}
}

func (h *healthService) Check(ctx context.Context) domain.HealthResponse {
	response := domain.HealthResponse{
		Status: "OK",
		Dependencies: []domain.Dependency{
			{
				Name:   "Database",
				Status: "OK",
			},
		},
	}

	err := h.database.Client().Ping(ctx, nil)
	if err != nil {
		response.Status = "ERROR"
		response.Dependencies[0].Status = "ERROR"
	}

	return response
}
