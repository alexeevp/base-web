package web

import (
	"github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"slices"
)

func middlewareAuth(e *echo.Echo) {
	e.Use(echojwt.WithConfig(echojwt.Config{
		Skipper: func(c echo.Context) bool {
			return slices.Contains([]string{"/login", "/auth", "/favicon.ico"}, c.Path())
		},
		SigningKey:  []byte(os.Getenv("JWT_SIGNKEY")),
		TokenLookup: "cookie:sid",
		ErrorHandler: func(c echo.Context, err error) error {
			return c.Redirect(http.StatusSeeOther, "/login")
		},
	}))
}
