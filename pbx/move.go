package pbx

import (
	"fmt"
	"go-getbox/utils"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
)

type Mover interface {
	MoveFileToFinished(project *Project) error
	Download(url string) (string, error)
}

type Move struct {
	request Requester
}

// MoveFileToFinished downloads and unzips render file
func (m *Move) MoveFileToFinished(project *Project) error {

	log.Printf("Start downloading file from %s", project.Render.URL)

	tempFile, err := m.Download(project.Render.URL)
	if err != nil {
		fmt.Printf("Couldn't download file %v.\n[ERROR] - %v\n", project.Render.URL, err)
		return err
	}

	log.Printf("Downloaded file %s\n", tempFile)

	finishedPath := m.getboxPath(m.request.GetConfig().GetString("finished_path") + project.Id)
	err = utils.Untar(tempFile, finishedPath)
	if err != nil {
		if err != io.EOF {
			log.Printf("Couldn't untar file %v.\n[ERROR] - %v\n", tempFile, err)
			return err
		}
	}

	log.Printf("Unpacked file %s and moved to finished %s\n", tempFile, finishedPath)

	return nil
}

// Download fetches file from url and returns temp local file path
func (m *Move) Download(url string) (string, error) {

	filename := m.extractFilename(url)
	downloadPath := m.getboxPath(m.request.GetConfig().GetString("temp_path"))
	err := m.request.DownloadFile(downloadPath+filename, url)
	return downloadPath + filename, err
}

func (m *Move) extractFilename(url string) string {
	return path.Base(url)
}

func (m *Move) getboxPath(directoryPath string) string {
	var getbox string
	getbox = m.request.GetConfig().GetString("getbox_path")

	if getbox == "" {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		getbox = filepath.Dir(dir) + "/"
	}

	path := getbox + directoryPath
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModeDir)
	}

	return path
}
