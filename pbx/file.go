package pbx

import (
	"fmt"
	"go-getbox/utils"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	jsonFilename = "info.json"
)

type File struct {
	request Requester
	config  *viper.Viper
}

// DownloadAndFinish downloads file, untars and moves to finished directory
func DownloadAndFinish(project *Project, request Requester, config *viper.Viper) error {

	file := &File{request, config}
	log.Printf("Start downloading file from %s", project.Render.URL)

	tempFile, err := file.download(project.Render.URL)
	if err != nil {
		fmt.Printf("Couldn't download file %v.\n[ERROR] - %v\n", project.Render.URL, err)
		return err
	}

	log.Printf("Downloaded file %s\n", tempFile)

	finishedPath := file.getboxPath(config.GetString("finished_path") + project.Id)
	err = utils.Untar(tempFile, finishedPath)
	if err != nil {
		if err != io.EOF {
			log.Printf("Couldn't untar file %v.\n[ERROR] - %v\n", tempFile, err)
			return err
		}
	}

	log.Printf("Unpacked file %s and moved to finished %s\n", tempFile, finishedPath)

	filePath := finishedPath + "/" + jsonFilename
	err = file.save(project, filePath)
	if err != nil {
		log.Printf("Couldn't write file %v", filePath)
	}

	log.Printf("Wrote file %v", filePath)

	return nil
}

// Download fetches file from url and returns temp local file path
func (f *File) download(url string) (string, error) {

	filename := f.extractFilename(url)
	downloadPath := f.getboxPath(f.config.GetString("temp_path"))
	err := f.request.DownloadFile(downloadPath+filename, url)
	return downloadPath + filename, err
}

func (f *File) extractFilename(url string) string {
	return path.Base(url)
}

func (f *File) getboxPath(directoryPath string) string {
	var getbox string
	getbox = f.config.GetString("getbox_path")

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

func (f *File) save(project *Project, filename string) error {
	err := ioutil.WriteFile(filename, project.RawJSON, 0644)
	return err
}
