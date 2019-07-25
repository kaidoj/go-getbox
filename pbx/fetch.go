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

type Fetcher interface {
	ProjectToSync() (Project, error)
	ProjectsToSync() bool
	FetchAndMoveProject(projects chan Project) error
}

type Fetch struct {
	request Requester
	move    Mover
	nrCores int
}

// NewFetcher Starts new fetcher instance
func NewFetcher(request Requester, move Mover) Fetcher {
	f := &Fetch{request, move, 0}
	return f
}

// ProjectsToSync syncs multiple projects concurrently
func (f *Fetch) ProjectsToSync() bool {
	f.nrCores = f.request.GetConfig().GetInt("nr_of_cores")

	projects := make(chan Project)
	var done []Project

	for i := 1; i <= f.nrCores; i++ {
		fmt.Printf("Run fetcher %d\n", i)
		wg.Add(1)
		go f.FetchAndMoveProject(projects)
	}

	wg.Add(1)
	go func() {
		for project := range projects {

			if project.Id != "" {
				fmt.Printf("Fetched project %s\n", project.Id)
			}

			done = append(done, project)

			// ALl the goroutines have finished
			// Close the channel
			if len(done) == f.nrCores {
				wg.Done()
			}
		}
	}()

	wg.Wait()
	close(projects)

	fmt.Println("No more renders to download")

	return true
}

// FetchAndMoveProject downloads projects, unzips and moves to finished directory
func (f *Fetch) FetchAndMoveProject(projects chan Project) error {
	defer wg.Done()

	project, err := f.ProjectToSync()
	defer f.appendProject(projects, project)
	if err != nil {
		return err
	}

	err = f.move.MoveFileToFinished(&project)
	if err != nil {
		return err
	}

	return nil
}

// ProjectToSync returns project that is ready for sync
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

func (f *Fetch) appendProject(projects chan Project, project Project) {
	projects <- project
}
