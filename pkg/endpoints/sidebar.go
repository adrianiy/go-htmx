package endpoints

import (
	"github.com/adrianiy/go-htmx/pkg/runner"
)

type Sidebar struct {
	Endpoints []runner.Endpoint
}

func NewSidebar() Sidebar {
	endpoints := runner.GetEndpoints()
	
	return Sidebar{
		Endpoints: endpoints,
	}
}

