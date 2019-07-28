package pbx

import (
	"fmt"
	"go-getbox/config"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	projectResponse = `{
		"id": "123",
		"name": "test",
		"render": {
			"url": "http://localhost"
		}	
	}`
)

func TestNext(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, projectResponse)
	}))
	defer ts.Close()

	config := config.Init("../tests")
	config.Set("host", extractHost(ts.URL))
	config.Set("schema", "http")
	r := NewRequest(config)
	_, err := Next(r)
	if err != nil {
		t.Errorf("Couldn't get next project. [ERROR] - %v", err)
	}
}

func TestSynced(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, projectResponse)
	}))
	defer ts.Close()

	config := config.Init("../tests")
	config.Set("host", extractHost(ts.URL))
	config.Set("schema", "http")
	r := NewRequest(config)

	p := &Project{}
	err := p.Synced(r)
	if err != nil {
		t.Errorf("Couldn't set status synced for project. [ERROR] - %v", err)
	}
}

func TestTryAgain(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, projectResponse)
	}))
	defer ts.Close()

	config := config.Init("../tests")
	config.Set("host", extractHost(ts.URL))
	config.Set("schema", "http")
	r := NewRequest(config)

	p := &Project{}
	err := p.TryAgain(r)
	if err != nil {
		t.Errorf("Couldn't set status synced for project. [ERROR] - %v", err)
	}
}
