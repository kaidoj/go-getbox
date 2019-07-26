package pbx

import (
	"go-getbox/config"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFileToFinished(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Handle("/", http.FileServer(http.Dir("../tests/files")))
	}))
	defer ts.Close()

	config := config.Init("../tests")
	r := &Request{}
	r.Config = config
	render := Render{0, ts.URL + "/render-file.tar"}
	p := &Project{}
	p.Id = "test123"
	p.Name = "filename"
	p.Render = render

	m := &Move{r}
	err := m.FileToFinished(p)

	if err != nil {
		t.Errorf("Couln't move file.[ERROR] - %v", err)
	}

}
