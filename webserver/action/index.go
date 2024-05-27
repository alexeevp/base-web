package action

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func Index(c echo.Context) error {
	return c.Render(http.StatusOK, "layout.html", map[string]interface{}{
		"name": "Dolly!",
	})
}
