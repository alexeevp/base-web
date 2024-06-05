package action

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func Index(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"title": "IndexProjectPage",
		"name":  "Index",
	})
}

func Login(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", nil)
}

func Auth(c echo.Context) error {

	cookie := new(http.Cookie)
	cookie.Name = "sid"
	cookie.Value = "valid-key"
	cookie.Expires = time.Now().Add(30 * 24 * time.Hour)
	c.SetCookie(cookie)
	return c.Redirect(http.StatusSeeOther, "/")
}

func Logout(c echo.Context) error {

	cookie := new(http.Cookie)
	cookie.Name = "sid"
	cookie.Value = "not-valid-key"
	cookie.Expires = time.Now().Add(-time.Hour)
	c.SetCookie(cookie)
	return c.Redirect(http.StatusSeeOther, "/login")
}
