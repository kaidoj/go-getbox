package pbx

import (
	"fmt"
	"log"
	"sync"

	"github.com/spf13/viper"
)

var (
	wg sync.WaitGroup
)

// Sync runs tasks
func Sync(config *viper.Viper) bool {
	nrCores := config.GetInt("nr_of_cores")
	request := NewRequest(config)

	projects := make(chan *Project)
	var done = 0

	for i := 1; i <= nrCores; i++ {
		fmt.Printf("Run fetcher %d\n", i)
		wg.Add(1)
		go Task(projects, request, config)
	}

	wg.Add(1)
	go func() {
		for project := range projects {
			if project != nil {
				fmt.Printf("Fetched project %s\n", project.Id)
			}
			done++
			// ALl the goroutines have finished
			// Close the channel and workgroup
			if done == nrCores {
				wg.Done()
			}
		}
	}()

	wg.Wait()
	close(projects)

	fmt.Println("No more renders to download")

	return true
}

//Task fetches next project and moves to finished
func Task(projects chan *Project, request Requester, config *viper.Viper) (*Project, error) {
	defer wg.Done()

	project, err := Next(request)
	defer appendProject(projects, project)
	if err != nil {
		return project, err
	}

	err = DownloadAndFinish(project, request, config)
	if err != nil {
		return project, err
	}

	err = project.Synced(request)
	if err != nil {
		log.Printf("Couldn't set synced status for project %v.\n[ERROR] - %s", project.Id, err)

		err = project.TryAgain(request)
		if err != nil {
			log.Printf("Couldn't set try again status for project %v.\n[ERROR] - %s", project.Id, err)
			return nil, err
		}

		return nil, err
	}

	log.Printf("Project %v status updated", project.Id)

	return project, nil
}

// appendProject send project through channel
func appendProject(projects chan *Project, project *Project) {
	projects <- project
}
