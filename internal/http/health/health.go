package health

import (
	"net/http"

	healthService "github.com/ADAGroupTcc/ms-channels-api/internal/services/health"
	"github.com/labstack/echo/v4"
)

type Health interface {
	Check(c echo.Context) error
}

type health struct {
	healthService healthService.HealthService
}

func New(healthService healthService.HealthService) Health {
	return &health{healthService}
}

func (h *health) Check(c echo.Context) error {
	ctx := c.Request().Context()

	res := h.healthService.Check(ctx)

	return c.JSON(http.StatusOK, res)
}
