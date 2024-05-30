package action

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func Second(c echo.Context) error {
	return c.Render(http.StatusOK, "second/second.html", map[string]interface{}{
		"name": "Inner page",
	})
}
