package server

import (
	"html/template"
	"io"
	"log"
	"github.com/labstack/echo/v4"
	"net/http"
	"github.com/adrianiy/go-htmx/pkg/endpoints"
)


type api struct {
	router http.Handler
}

type Server interface {
	Router() http.Handler
}

type TemplateRenderer struct {
	templates *template.Template
}

type Page struct {
	Title string
	Sidebar endpoints.Sidebar
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func New() Server {
	Page := Page{
		Title: "Go htmx",
		Sidebar: endpoints.newSidebar(),
	}
	tmpl, err := template.ParseGlob("./templates/*.html")

	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}

	a := &api{}
	
	r := echo.New()

	r.Renderer = &TemplateRenderer{
		templates: tmpl,
	}

	r.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "main", Page)
	})

	r.GET("/endpoints", func(c echo.Context) error {
		return c.Render(http.StatusOK, "endpoints", nil)
	})
	
	a.router = r
	
	return a
}

func (a *api) Router() http.Handler {
	return a.router
}
