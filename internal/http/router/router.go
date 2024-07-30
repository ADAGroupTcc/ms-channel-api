package router

import (
	"github.com/ADAGroupTcc/ms-channels-api/config"
	"github.com/ADAGroupTcc/ms-channels-api/internal/http/middlewares"
	"github.com/labstack/echo/v4"
)

func SetupRouter(dependencies *config.Dependencies) *echo.Echo {
	e := echo.New()
	e.Use(middlewares.ErrorIntercepter())

	e.GET("/health", dependencies.HealthHandler.Check)

	v1 := e.Group("/v1")
	v1.POST("/channels", dependencies.Handler.Create)
	v1.GET("/channels/:id", dependencies.Handler.Get)
	v1.GET("/channels", dependencies.Handler.List)
	v1.PATCH("/channels/:id", dependencies.Handler.Update)
	v1.DELETE("/channels/:id", dependencies.Handler.Delete)

	return e
}
