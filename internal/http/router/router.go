package router

import (
	"github.com/ADAGroupTcc/ms-channels-api/config"
	"github.com/ADAGroupTcc/ms-channels-api/internal/http/middlewares"
	"github.com/labstack/echo/v4"
)

func SetupRouter(dependencies *config.Dependencies) *echo.Echo {
	e := echo.New()
	e.Use(middlewares.ErrorIntercepter())

	e.GET("/health", dependencies.HealthHandler.Check, middlewares.ErrorIntercepter())

	v1 := e.Group("/v1")
	v1.POST("/channels", dependencies.Handler.Create, middlewares.ErrorIntercepter())
	v1.GET("/channels/:id", dependencies.Handler.Get, middlewares.ErrorIntercepter())
	v1.GET("/channels", dependencies.Handler.List, middlewares.ErrorIntercepter())
	v1.PATCH("/channels/:id", dependencies.Handler.Update, middlewares.ErrorIntercepter())
	v1.DELETE("/channels/:id", dependencies.Handler.Delete, middlewares.ErrorIntercepter())

	return e
}
