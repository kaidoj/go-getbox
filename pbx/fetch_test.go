package pbx

import (
	"fmt"
	"go-getbox/config"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	response = `{
		"id": "123",
		"name": "test",
		"render": {
			"url": "http://localhost"
		}	
	}`
)

func TestProjectsToSync(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, response)
	}))
	defer ts.Close()

	config := config.Init("../tests")
	config.Set("host", extractHost(ts.URL))
	config.Set("Schema", "http")

	getbox := &Getbox{}
	getbox.Config = config
	r := &Request{}
	r.Config = config
	m := &Move{r}
	fetch := NewFetcher(r, m)

	res := fetch.ProjectsToSync()
	if !res {
		t.Errorf("Something wrong")
	}
}
