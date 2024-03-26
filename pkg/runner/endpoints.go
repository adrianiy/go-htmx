package runner

import (
	"log"
	"io/ioutil"
	"context"
	"github.com/google/go-github/github"
	"gopkg.in/yaml.v3"
)

const (
	owner    = "adrianiy"
	repo     = "go-htmx"
	basePath = "/knowledge-base/endpoints"
)

var client = github.NewClient(nil)

type Node struct {
	Node string
	Sql string
}

type Endpoint struct {
	Id int
	Path string
	Name string
	Description string
	Nodes []Node
}

var id int = 0

func newEndpoint(path, name, description string) Endpoint {
	return Endpoint{Id: id, Path: path, Name: name, Description: description}
}

func getEndpoints() []Endpoint {
	_, dirContent, _, err := client.Repositories.GetContents(context.Background(), owner, repo, basePath, nil)

	if err != nil {
		log.Fatalf("Error getting contents: %v", err)
	}

	var endpoints []Endpoint
	
	for _, file := range dirContent {
		buf, err := ioutil.ReadFile(file.GetDownloadURL())

		if err != nil {
			log.Fatalf("Error getting file content: %v", err)
		}

		e := &Endpoint{}

		err = yaml.Unmarshal(buf, e)

		if err != nil {
			log.Fatalf("Error unmarshalling yaml: %v", err)
		}
		
		id++
		e.Id = id
		
		endpoints = append(endpoints, *e)
	}
		
	return endpoints
}


