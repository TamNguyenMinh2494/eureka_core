package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Permit(roles ...string) func(echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			profile := c.Get("user").(map[string]interface{})
			for _, role := range roles {
				if profile["role"] == role {
					if err := next(c); err != nil {
						c.Error(err)
					}
					return nil
				}
			}
			return echo.NewHTTPError(http.StatusMethodNotAllowed, map[string]interface{}{"message": "Permission denied"})
		}
	}
}

func Deny(roles ...string) func(echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			profile := c.Get("user").(map[string]interface{})
			for _, role := range roles {
				if profile["role"] == role {
					return echo.NewHTTPError(http.StatusMethodNotAllowed, map[string]interface{}{"message": "Permission denied"})

				}
			}
			if err := next(c); err != nil {
				c.Error(err)
			}
			return nil
		}
	}
}
