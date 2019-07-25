package config

import (
	"testing"
)

func TestInit(t *testing.T) {
	config := Init("../tests")
	if config == nil {
		t.Errorf("No config found")
	}

	v := "localhost"
	res := config.GetString("host")
	if res != v {
		t.Errorf("want %v; got %v", v, res)
	}
}
