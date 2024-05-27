package webserver

import (
	"bms/webserver/action"
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
	e.GET("/", action.Index)
	e.Static("/static", "static")
	e.HideBanner = true
	e.Renderer = &TemplateRenderer{
		templates: template.Must(template.ParseGlob("webserver/templates/*.html")),
	}
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

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}
