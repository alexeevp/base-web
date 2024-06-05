package web

import (
	"bms/web/action"
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"
)

var e *echo.Echo

func Run() {

	e = echo.New()
	//e.Use(middleware.Logger())
	e.Use(middlewareAuth())
	e.GET("/", action.Index)
	e.GET("/favicon.ico", func(c echo.Context) error {
		return c.String(http.StatusNoContent, "")
	})
	e.GET("/login", action.Login)
	e.POST("/auth", action.Auth)
	e.POST("/logout", action.Logout)
	e.GET("/second", action.Second)
	e.Static("/static", "web/static")
	e.HideBanner = true
	e.Renderer = &TemplateRenderer{}

	go func() {
		err := e.Start(":1323")
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatal(err)
		}
	}()
}

func Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := e.Shutdown(ctx)
	if err != nil {
		e.Logger.Fatal(err)
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

	t.templates = template.Must(template.ParseFiles("web/templates/"+name, "web/templates/layout.html"))
	return t.templates.ExecuteTemplate(w, "layout.html", data)
}
