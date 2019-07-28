package pbx

import (
	"fmt"
	"go-getbox/config"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestDownloadAndFinish(t *testing.T) {
	file := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Handle("/files/", http.FileServer(http.Dir("../tests/files")))
	}))
	defer file.Close()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, projectResp(file.URL, r.Method))
	}))
	defer ts.Close()

	config := config.Init("../tests")
	config.Set("host", extractHost(ts.URL))
	config.Set("schema", "http")
	r := NewRequest(config)
	project := &Project{}
	project.Id = "123"
	project.Name = "Test"
	project.Render = Render{0, file.URL + "/render-file.tar"}
	err := DownloadAndFinish(project, r, config)
	if err != nil {
		t.Errorf("Download and finish failed")
	}
}

func TestSave(t *testing.T) {
	config := config.Init("../tests")
	project := &Project{}
	project.Id = "123"
	project.Name = "Test"
	project.RawJSON = []byte("{testjson}")

	file := &File{}
	file.config = config
	finishedPath := config.GetString("finished_path") + project.Id + "/"
	filePath := file.getboxPath(finishedPath) + jsonFilename
	err := file.save(project, filePath)
	if err != nil {
		t.Errorf("Couldn't write file %v", filePath)
	}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("File %v not found", filePath)
	}
}
