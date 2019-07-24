package pbx

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
)

var (
	wg sync.WaitGroup
)

type Fetcher interface {
	ProjectToSync() Project
	ProjectsToSync()
	FetchAndMoveProject(projects chan Project, done chan bool) error
}

type Fetch struct {
	request Requester
	move    Mover
	nrCores int
}

//ProjectsToSync syncs multiple projects concurrently
func (f *Fetch) ProjectsToSync() {
	f.nrCores = f.request.GetConfig().GetInt("nr_of_cores")

	projects := make(chan Project)

	for i := 1; i <= f.nrCores; i++ {
		log.Printf("Run fetcher %d", i)
		wg.Add(1)
		go f.FetchAndMoveProject(projects)
	}

	go func() {
		for project := range projects {
			log.Printf("Fetched project %s", project.Id)
		}
	}()

	wg.Wait()

	fmt.Println("No more renders to download")
	os.Exit(1)
}

//FetchAndMoveProject downloads projects, unzips and moves to finished directory
func (f *Fetch) FetchAndMoveProject(projects chan Project) error {
	project, err := f.ProjectToSync()

	if err != nil {
		wg.Done()
		return err
	}

	f.move.MoveFileToFinished(&project)
	projects <- project

	if len(projects) == f.nrCores {
		wg.Done()
	}

	return nil
}

//ProjectToSync returns project that is ready for sync
func (f *Fetch) ProjectToSync() (Project, error) {
	endpoint := "projects/next"
	res, err := f.request.Get(endpoint)

	if err != nil {
		log.Printf("No results returned for request %v.\n[ERROR] - %s", endpoint, err)
		return Project{}, err
	}

	var project Project
	err = json.Unmarshal(res, &project)

	if err != nil {
		return project, err
	}

	return project, nil
}
