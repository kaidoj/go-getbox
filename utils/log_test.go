package utils

import (
	"fmt"
	"go-getbox/config"
	"log"
	"os"
	"strings"
	"testing"
)

func TestLogOutputWithConfigVar(t *testing.T) {
	config := config.Init("../tests")
	config.Set("logs_dir", "../tests/logs/")
	file := LogOutput(config)
	defer file.Close()
	log.SetOutput(file)
	v := "Test error to log"
	log.Println(v)
	res := readLastLine(file.Name())
	if !strings.Contains(res, v) {
		t.Errorf("want %v; got %v", v, res)
	}
}

func readLastLine(fname string) string {
	file, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		fmt.Println(err)
	}

	buf := make([]byte, 32)
	n, err := file.ReadAt(buf, fi.Size()-int64(len(buf)))
	if err != nil {
		fmt.Println(err)
	}
	buf = buf[:n]
	fmt.Printf("%s", buf)

	return string(buf)
}
