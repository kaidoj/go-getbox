package utils

import (
	"os"
	"testing"
)

func TestUntar(t *testing.T) {
	dest := "../tests/getbox/finished/render-file"
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		os.MkdirAll(dest, os.ModeDir)
	}

	err := Untar("../tests/files/render-file.tar", dest)
	if err != nil {
		t.Errorf("%v", err)
	}

	v := dest + "/page_0000.jpg"
	if _, err := os.Stat(v); os.IsNotExist(err) {
		t.Errorf("File %v doesn't exist.\n[ERROR] - %v", v, err)
	}

}
