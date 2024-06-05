package web

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"slices"
)

func middlewareAuth() echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		Skipper: func(c echo.Context) bool {
			return slices.Contains([]string{"/login", "/auth", "/favicon.ico"}, c.Path())
		},
		KeyLookup: "cookie:sid",
		Validator: func(key string, c echo.Context) (bool, error) {
			return key == "valid-key", nil
		},
		ErrorHandler: func(err error, c echo.Context) error {
			return c.Redirect(http.StatusSeeOther, "/login")
		},
	})
}
