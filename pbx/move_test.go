package pbx

import (
	"getbox/config"
	"testing"
)

func TestMoveFileToFinished(t *testing.T) {
	config := config.Init("../tests")
	r := &Request{}
	r.Config = config
	render := Render{0, "http://localhost"}
	p := &Project{}
	p.Id = "test123"
	p.Name = "filename"
	p.Render = render

	m := &Move{r}
	_, err := m.MoveFileToFinished(p)

	if err != nil {
		t.Errorf("Couln't move file. [ERROR] - %v", err)
	}

}
