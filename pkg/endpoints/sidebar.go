package endpoints

import (
	"github.com/adrianiy/go-htmx/pkg/runner"
)

type Sidebar struct {
	Endpoints []runner.Endpoint
}

func newSidebar() Sidebar {
	endpoints := runner.getEndpoints()
	
	return Sidebar{
		Endpoints: endpoints,
	}
}

