package pbx

import (
	"fmt"
	"go-getbox/config"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSync(t *testing.T) {
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
	res := Sync(config)

	if !res {
		t.Errorf("Couldn't sync")
	}
}

func TestTask(t *testing.T) {
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
	projects := make(chan *Project)
	wg.Add(1)
	go Task(projects, r, config)

	for project := range projects {
		if project == nil {
			t.Errorf("Couldn't run task.")
		}

		if project.Id != "123" {
			t.Errorf("Project id not found.")
		}

		close(projects)
	}
}

func projectResp(customURL, method string) string {
	if method == http.MethodPost {
		return `{
			"id": "123",
			"name": "test",
			"render": {
				"url": ""
			}
		}`
	}

	return `{
		"id": "123",
		"name": "test",
		"render": {
			"url": "` + customURL + `/files/render-file.tar"
		}
	}`
}
