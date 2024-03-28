package components

import (
	"github.com/adrianiy/go-htmx/pkg/entities"
)

type Sidebar struct {
	Endpoints []entities.Endpoint
}

func NewSidebar() Sidebar {
	endpoints := entities.GetEndpoints()
	
	return Sidebar{
		Endpoints: endpoints,
	}
}

