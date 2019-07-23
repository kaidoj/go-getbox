package pbx

import (
	"encoding/json"
	"fmt"
	"log"
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
	nrCores := f.request.Config.GetInt("nr_of_cores")

	in := make(chan Project)
	out := make(chan Project)

	projects := []Project

	for i := 1; i <= nrCores; i++ {
		fmt.Printf("ProjectToSync %d", i)
		project := f.ProjectToSync()
		go f.FetchAndMoveProject(project, syncing)
	}

	for o := range out {
		fmt.Println("called out")
		categoriesList = append(categoriesList, o)

		if lenCategories == len(categoriesList) {
			close(out)
		}
	}
}

func (f *Fetcher) FetchAndMoveProject(project *Project, syncing chan) {

	project := ProjectToSync()
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
		log.Fatalf("Couldn't parse json in request %v.[ERROR] - %s", endpoint, err)
	}

	return project
}
