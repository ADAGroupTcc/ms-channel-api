package middlewares

import (
	"github.com/ADAGroupTcc/ms-channels-api/exceptions"
	"github.com/labstack/echo/v4"
)

func ErrorIntercepter() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err != nil {
				res := exceptions.HandleExceptions(err)
				return c.JSON(res.Code, res)
			}
			return err
		}
	}
}
