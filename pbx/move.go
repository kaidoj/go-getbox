package pbx

import (
	"fmt"
	"getbox/utils"
	"log"
	"os"
	"path"
	"path/filepath"
)

type Mover interface {
	MoveFileToFinished(project *Project) ([]string, error)
	Download(url string) (string, error)
}

type Move struct {
	request Requester
}

//MoveFileToFinished downloads and unzips render file
func (m *Move) MoveFileToFinished(project *Project) ([]string, error) {

	log.Printf("Start downloading file from %s", project.Render.URL)

	tempFile, err := m.Download(project.Render.URL)
	if err != nil {
		fmt.Printf("Couldn't download file %v.\n[ERROR] - %v", project.Render.URL, err)
		return nil, err
	}

	log.Printf("Downloaded file %s", tempFile)

	finishedPath := m.getboxPath(m.request.GetConfig().GetString("finished_path") + project.Id)
	files, err := utils.Unzip(tempFile, finishedPath)
	if err != nil {
		log.Printf("Couldn't unzip file %v.\n[ERROR] - %v", tempFile, err)
		return nil, err
	}

	log.Printf("Unzipped file %s and moved to finished %s", tempFile, finishedPath)

	return files, nil
}

//Download fetches file from url and returns temp local file path
func (m *Move) Download(url string) (string, error) {

	filename := m.extractFilename(url)
	downloadPath := m.getboxPath(m.request.GetConfig().GetString("temp_path"))
	err := m.request.DownloadFile(downloadPath, url)
	return downloadPath + filename, err
}

func (m *Move) extractFilename(url string) string {
	return path.Base(url)
}

func (m *Move) getboxPath(directoryPath string) string {
	var getbox string
	getbox = m.request.GetConfig().GetString("getbox_path")
	newpath := filepath.Join(".", "public")
	fmt.Println(newpath)

	if getbox == "" {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		getbox = filepath.Dir(dir) + "/"
	}

	path := getbox + directoryPath
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModeDir)
	}

	return path
}
