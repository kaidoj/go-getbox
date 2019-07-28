package pbx

import (
	"encoding/json"
	"log"
)

const (
	syncStatusSuccess  = "SUCCESS"
	syncStatusTryAgain = "TRY_AGAIN"
)

type Project struct {
	Id     string
	Name   string
	Render Render
}

type Render struct {
	Status int
	URL    string
}

// Next fetches next project from api
func Next(request Requester) (*Project, error) {
	endpoint := "projects/next"
	res, err := request.Get(endpoint)

	if err != nil {
		log.Printf("No results returned for request %v.\n[ERROR] - %s", endpoint, err)
		return &Project{}, err
	}

	var project *Project
	err = json.Unmarshal(res, &project)
	if err != nil {
		return project, err
	}

	return project, nil
}

// Synced sets getbox status to synced
func (project *Project) Synced(request Requester) error {

	payload := map[string]string{
		"status": syncStatusSuccess,
	}

	_, err := request.Post("projects/"+project.Id, payload)
	if err != nil {
		log.Printf("Couldn't set status %v for project %v.\n[ERROR] - %s", syncStatusSuccess, project.Id, err)
		return err
	}

	return nil
}

// TryAgain sets getbox status to try again
func (project *Project) TryAgain(request Requester) error {

	payload := map[string]string{
		"status": syncStatusTryAgain,
	}

	_, err := request.Post("projects/"+project.Id, payload)
	if err != nil {
		log.Printf("Couldn't set status %v for project %v.\n[ERROR] - %s", syncStatusTryAgain, project.Id, err)
		return err
	}

	return nil
}
