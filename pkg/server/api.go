package server

import (
	"html/template"
	"io"
	"log"
	"github.com/labstack/echo/v4"
	"net/http"
	"github.com/adrianiy/go-htmx/pkg/components"
	"github.com/adrianiy/go-htmx/pkg/entities"
	"github.com/google/uuid"
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

type Content struct {
	Content components.Content
}

type Page struct {
	Title string
	Sidebar components.Sidebar
	Content components.Content
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func New() Server {
	Sidebar := components.NewSidebar()
	Page := Page{
		Title: "Go htmx",
		Sidebar: Sidebar,
		Content: components.Content{},
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

	r.Static("/", "assets")

	r.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "main", Page)
	})

	r.GET("/endpoint/:id", func(c echo.Context) error {
		idStr := c.Param("id")

		id, err := uuid.Parse(idStr)

		if err != nil {
			log.Printf("Error converting id to uuid: %v", err)
			
			return c.String(http.StatusBadRequest, "Invalid id")
		}

		endpoint, err1 := entities.FindEndpoint(Sidebar.Endpoints, id)

		if err1 != nil {
			log.Printf("Error finding endpoint: %v", err1)

			// redirect to home
			return c.Redirect(http.StatusFound, "/")
		}

		Page.Content = components.NewContent(endpoint)

		err2 := c.Render(http.StatusOK, "main", Page)

		if err2 != nil {
			log.Printf("Error rendering main: %v", err1)

			return c.String(http.StatusInternalServerError, "Internal server error")
		}

		return nil
	})
		

	r.GET("/endpoints", func(c echo.Context) error {
		return c.Render(http.StatusOK, "endpoints", nil)
	})

	r.GET("/content/:id", func(c echo.Context) error {
		idStr := c.Param("id")

		id, err := uuid.Parse(idStr)

		if err != nil {
			log.Printf("Error converting id to uuid: %v", err)

			return c.String(http.StatusBadRequest, "Invalid id")
		}
		
		endpoint, err1 := entities.FindEndpoint(Sidebar.Endpoints, id)

		if err1 != nil {
			log.Printf("Error finding endpoint: %v", err1)

			return c.String(http.StatusNotFound, "Endpoint not found")
		}
		
		Content := Content{
			Content: components.NewContent(endpoint),
		}
		
		err2 := c.Render(http.StatusOK, "content", Content)

		if err2 != nil {
			log.Fatalf("Error rendering content: %v", err1)

			return c.String(http.StatusInternalServerError, "Internal server error")
		}

		return nil
	})
	
	a.router = r
	
	return a
}

func (a *api) Router() http.Handler {
	return a.router
}
