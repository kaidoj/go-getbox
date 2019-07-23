package pbx

import (
	"getbox/config"
	"testing"
)

func TestIsNumberOfCoresSet(t *testing.T) {
	config := config.Init("../tests")
	gbox := &Getbox{config}
	err := gbox.isNrCoresSet()
	if err == false {
		t.Errorf("Number of cores not set")
	}
}
