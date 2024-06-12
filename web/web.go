package web

import (
	"bms/web/action"
	"context"
	"errors"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"log"
	"net/http"
	"slices"
	"time"
)

type Webserver struct {
	*echo.Echo
	jwtSignKey string
}

func NewWebserver(jwtSignKey string) *Webserver {
	w := &Webserver{
		echo.New(),
		jwtSignKey,
	}
	w.DefaultConfig()
	return w
}

func (w *Webserver) DefaultConfig() {
	w.Use(echojwt.WithConfig(echojwt.Config{
		Skipper: func(c echo.Context) bool {
			return slices.Contains([]string{"/login", "/auth", "/favicon.ico"}, c.Path())
		},
		SigningKey:  []byte(w.jwtSignKey),
		TokenLookup: "cookie:sid",
		ErrorHandler: func(c echo.Context, err error) error {
			return c.Redirect(http.StatusSeeOther, "/login")
		},
	}))
	//w.Use(middleware.Logger())

	w.GET("/", action.Index)
	w.Static("/static", "web/static")
	w.GET("/favicon.ico", func(c echo.Context) error {
		return c.String(http.StatusNoContent, "")
	})
	w.GET("/login", Login)
	w.POST("/auth", Auth)
	w.POST("/logout", Logout)
	w.GET("/second", action.Second)

	w.HideBanner = true
	w.Renderer = &TemplateRenderer{}
}

func (w *Webserver) Run() {
	go func() {
		err := w.Start(":1323")
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			w.Logger.Fatal(err)
		}
	}()
}

func (w *Webserver) Off() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := w.Shutdown(ctx)
	if err != nil {
		w.Logger.Fatal(err)
	}
	log.Println("Webserver shutdown complete")
}

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	if c.Path() == "/login" {
		t.templates = template.Must(template.ParseFiles("web/templates/" + name))
		return t.templates.ExecuteTemplate(w, name, data)
	}

	viewContext, isMap := data.(map[string]interface{})
	if !isMap {
		return errors.New("Can't render without map data.")
	}
	viewContext["reverse"] = c.Echo().Reverse
	viewContext["user"] = "u" //c.Get("user").(*jwt.Token).Claims.(*jwtCustomClaims)

	t.templates = template.Must(template.ParseFiles("web/templates/"+name, "web/templates/layout.html"))
	return t.templates.ExecuteTemplate(w, "layout.html", data)
}
