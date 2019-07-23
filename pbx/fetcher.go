package pbx

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

var (
	wg sync.WaitGroup
)

type Fetch interface {
	ProjectToSync() Project
	ProjectsToSync()
}

type Fetcher struct {
	request *Request
}

type Project struct {
	Id     string
	Name   string
	Render Render
}

type Render struct {
	Status int
	URL    string
}

func (f *Fetcher) ProjectsToSync() {
	defer wg.Done()
	nrCores := f.request.Config.GetInt("nr_of_cores")

	projects := make(chan Project)

	for i := 1; i <= nrCores; i++ {
		fmt.Printf("Run fetcher %d", i)
		wg.Add(1)
		go f.FetchProject(projects)
	}

	for project := range projects {
		fmt.Printf("Fetched project %s", project.Name)
	}
}

func (f *Fetcher) FetchProject(projects chan Project) {
	project := f.ProjectToSync()
	projects <- project
}

//ProjectsToSync returns projects that are ready for sync
func (f *Fetcher) ProjectToSync() Project {
	endpoint := "projects/next"
	res, err := f.request.Get(endpoint)

	if err != nil {
		log.Fatalf("No results returned for request %v.[ERROR] - %s", endpoint, err)
	}

	var project Project
	err = json.Unmarshal(res, &project)

	if err != nil {
		fmt.Println("No more renders to download")
		log.Fatalf("Couldn't parse json in request %v.[ERROR] - %s", endpoint, err)
	}

	return project
}
