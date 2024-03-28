package entities

import (
	"crypto/sha1"
	"encoding/hex"
	"os"
	"strconv"
	"strings"
	"log"
	"fmt"
	"errors"
	"path/filepath"
	"io/ioutil"
	"context"
	"github.com/google/go-github/github"
	"gopkg.in/yaml.v3"
	"github.com/google/uuid"
)

const (
	owner    = "adrianiy"
	repo     = "go-htmx"
	basePath = "knowledge-base/endpoints"
)

var client = github.NewClient(nil)

type Node struct {
	Node string
	Sql string
}

type Endpoint struct {
	Id uuid.UUID
	Path string
	Name string
	Description string
	Nodes []Node
}

type Endpoints []Endpoint

func NewEndpoint(path, name, description string, nodes []Node) Endpoint {
	id, err := uuid.NewV6()
	fmt.Println("id:", id)

	if err != nil {
		log.Fatalf("Error generating uuid: %v", err)
	}
	
	return Endpoint{Id: id, Path: path, Name: strings.Trim(name), Description: description, Nodes: nodes}
}

func FindEndpoint(endpoints []Endpoint, id uuid.UUID) (Endpoint, error) {
	for _, e := range endpoints {
		if e.Id == id {
			return e, nil
		}
	}

	return Endpoint{}, errors.New("Endpoint not found")
}

func (e *Endpoint) AddNode(node, sql string) {
	n := Node{Node: node, Sql: sql}
	e.Nodes = append(e.Nodes, n)
}

func (e *Endpoint) AddDescription(description string) {
	e.Description = description
}

func (e *Endpoint) AddName(name string) {
	e.Name = name
}

func (e *Endpoint) AddPath(path string) {
	e.Path = path
}

func  getContents(path string) []Endpoint {
	endpoints := []Endpoint{}
	
	_, directoryContent, _, err := client.Repositories.GetContents(context.Background(), owner, repo, path, nil)

	if err != nil {
		log.Fatalf("Error getting contents: %v", err)
	}
	
	for _, c := range directoryContent {
		fmt.Println(*c.Type, *c.Path, *c.Size, *c.SHA)

		dir, _ := os.Getwd()
		local := filepath.Join(dir, *c.Path)
		fmt.Println("local:", local)

		switch *c.Type {
		case "file":
			_, err := os.Stat(local)
			
			if err == nil {
				b, err1 := ioutil.ReadFile(local)
				
				if err1 == nil {
					sha := calculateGitSHA1(b)
					
					if *c.SHA == hex.EncodeToString(sha) {
						endpoints = addEndpointFromLocalFile(endpoints, local)
						continue
					}
				}
			}
			downloadContents(c, local)
			endpoints = addEndpointFromLocalFile(endpoints, local)
		case "dir":
			getContents(filepath.Join(path, *c.Path))
		}
	}

	return endpoints
}

func GetEndpoints() []Endpoint {
	return getContents(basePath)
}

func downloadContents(file *github.RepositoryContent, local string) {
	rc, err := client.Repositories.DownloadContents(context.Background(), owner, repo, *file.Path, nil)

	if err != nil {
		log.Fatalf("Error downloading contents: %v", err)
	}

	buf, err := ioutil.ReadAll(rc)

	if err != nil {
		log.Fatalf("Error reading content: %v", err)
	}

	err = ioutil.WriteFile(local, buf, 0644)

	if err != nil {
		log.Fatalf("Error writing content: %v", err)
	}
}

func addEndpointFromLocalFile(endpoints []Endpoint, path string) []Endpoint {
	fmt.Println("path:", path)
	b, err := ioutil.ReadFile(path)

	fmt.Println("b:", string(b))

	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	e := Endpoint{}

	err = yaml.Unmarshal(b, &e)

	fmt.Println("e:", e)

	if err != nil {
		log.Fatalf("Error unmarshalling yaml: %v", err)
	}

	return append(endpoints, NewEndpoint(e.Path, e.Name, e.Description, e.Nodes)
}

// calculateGitSHA1 computes the github sha1 from a slice of bytes.
// The bytes are prepended with: "blob " + filesize + "\0" before runing through sha1.
func calculateGitSHA1(contents []byte) []byte {
	contentLen := len(contents)
	blobSlice := []byte("blob " + strconv.Itoa(contentLen))
	blobSlice = append(blobSlice, '\x00')
	blobSlice = append(blobSlice, contents...)
	h := sha1.New()
	h.Write(blobSlice)
	bs := h.Sum(nil)
	return bs
}
