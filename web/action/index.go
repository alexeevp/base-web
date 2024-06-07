package action

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"time"
)

const tokenLifetime = 72

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

	username := c.FormValue("username")
	if username != "adminuser" {
		return echo.ErrUnauthorized
	}

	signedToken, err := buildToken(username)
	if err != nil {
		return err
	}

	cookie := new(http.Cookie)
	cookie.Name = "sid"
	cookie.Value = signedToken
	cookie.Expires = time.Now().Add(time.Hour * tokenLifetime)
	c.SetCookie(cookie)
	return c.Redirect(http.StatusSeeOther, "/")
}

func Logout(c echo.Context) error {

	cookie := new(http.Cookie)
	cookie.Name = "sid"
	cookie.Value = "nvk"
	cookie.Expires = time.Now().Add(-time.Hour)
	c.SetCookie(cookie)
	return c.Redirect(http.StatusSeeOther, "/login")
}

type jwtCustomClaims struct {
	Username string `json:"username"`
	Admin    bool   `json:"admin"`
	jwt.RegisteredClaims
}

func buildToken(username string) (string, error) {

	claims := &jwtCustomClaims{
		username,
		true,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * tokenLifetime)),
		},
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := rawToken.SignedString([]byte(os.Getenv("JWT_SIGNKEY")))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
