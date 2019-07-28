package pbx

import (
	"go-getbox/config"
	"testing"
)

func TestIsNumberOfCoresSet(t *testing.T) {
	config := config.Init("../tests")
	err := IsNrCoresSet(config)
	if err == false {
		t.Errorf("Number of cores not set")
	}
}
