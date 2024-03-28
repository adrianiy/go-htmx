package components

import (
	"github.com/adrianiy/go-htmx/pkg/entities"
)

type Content entities.Endpoint

func NewContent(endpoint entities.Endpoint) Content {
	return Content(endpoint)
}

